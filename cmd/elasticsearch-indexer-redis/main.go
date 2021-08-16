package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rbac/cmd/internal"
	internaldomain "rbac/internal"
	"rbac/internal/elasticsearch"
	"rbac/internal/envvar"
	"rbac/internal/memcached"
	"strings"
	"syscall"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	errC, err := run(env)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(env string) (<-chan error, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction %w", err)
	}

	if err := envvar.Load(env); err != nil {
		return nil, fmt.Errorf("envvar.Load %w", err)
	}

	vault, err := internal.NewVaultProvider()
	if err != nil {
		return nil, fmt.Errorf("internal.NewVaultProvider %w", err)
	}

	conf := envvar.New(vault)

	//-

	rdb, err := internal.NewRedis(conf)
	if err != nil {
		return nil, fmt.Errorf("newRedis %w", err)
	}

	//-

	_, err = internal.NewOTExporter(conf)
	if err != nil {
		return nil, fmt.Errorf("newOTExporter %w", err)
	}

	//-

	es, err := internal.NewElasticSearch(conf)
	if err != nil {
		return nil, fmt.Errorf("internal.NewElasticSearch %w", err)
	}

	mem, err := internal.NewMemcached(conf)
	if err != nil {
		return nil, fmt.Errorf("internal.NewMemcached %w", err)
	}
	search := elasticsearch.NewRBAC(es)
	mclient := memcached.NewRBAC(mem, search, logger)

	srv := &Server{
		logger: logger,
		rdb:    rdb,
		cache:  mclient,
		done:   make(chan struct{}),
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			logger.Sync()
			rdb.Close()
			stop()
			cancel()
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving")

		if err := srv.ListenAndServe(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}

type Server struct {
	logger *zap.Logger
	rdb    *redis.Client
	pubsub *redis.PubSub
	cache  *memcached.RBAC
	done   chan struct{}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	pubsub := s.rdb.PSubscribe(context.Background(), "accounts.*")

	_, err := pubsub.Receive(context.Background())
	if err != nil {
		return fmt.Errorf("pubsub.Receive %w", err)
	}

	s.pubsub = pubsub

	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			s.logger.Info(fmt.Sprintf("Received message: %s", msg.Channel))

			switch msg.Channel {
			case "accounts.event.created":
				var account internaldomain.Account
				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&account); err != nil {
					s.logger.Info("Ignoring message, invalid", zap.Error(err))
					continue
				}
				if err := s.cache.IndexAccount(context.Background(), account); err != nil {
					s.logger.Info("Couldn't index account", zap.Error(err))
				}
			case "accounts.event.updated":
				var profile internaldomain.Profile
				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&profile); err != nil {
					s.logger.Info("Ignoring message, invalid", zap.Error(err))
					continue
				}
				if err := s.cache.DeleteProfile(context.Background(), profile.Id); err != nil {
					s.logger.Info("Couldn't delete profile", zap.Error(err))
				}
				if err := s.cache.IndexProfile(context.Background(), profile); err != nil {
					s.logger.Info("Couldn't index profile", zap.Error(err))
				}

			case "accounts.event.deleted":
				var id string
				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&id); err != nil {
					s.logger.Info("Ignoring message, invalid", zap.Error(err))
					continue
				}
				if err := s.cache.DeleteProfile(context.Background(), id); err != nil {
					s.logger.Info("Couldn't delete task", zap.Error(err))
				}
			}
		}

		s.logger.Info("No more messages to consume. Exiting.")

		s.done <- struct{}{}
	}()

	return nil
}

// Shutdown ...
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")

	s.pubsub.Close()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context.Done: %w", ctx.Err())

		case <-s.done:
			return nil
		}
	}
}
