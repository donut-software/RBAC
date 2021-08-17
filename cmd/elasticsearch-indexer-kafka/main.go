package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"rbac/cmd/events"
	"rbac/cmd/internal"
	internaldomain "rbac/internal"
	"rbac/internal/elasticsearch"
	"rbac/internal/envvar"
	"rbac/internal/memcached"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	es, err := internal.NewElasticSearch(conf)
	if err != nil {
		return nil, fmt.Errorf("internal.NewElasticSearch %w", err)
	}

	kafka, err := internal.NewKafkaConsumer(conf, "elasticsearch-indexer")
	if err != nil {
		return nil, fmt.Errorf("internal.NewKafkaConsumer %w", err)
	}

	//-

	_, err = internal.NewOTExporter(conf)
	if err != nil {
		return nil, fmt.Errorf("newOTExporter %w", err)
	}

	//-
	mem, err := internal.NewMemcached(conf)
	if err != nil {
		return nil, fmt.Errorf("internal.NewMemcached %w", err)
	}
	search := elasticsearch.NewRBAC(es, 100)
	mClient := memcached.NewRBAC(mem, search, logger)

	srv := &Server{
		logger: logger,
		kafka:  kafka,
		rbac:   mClient,
		doneC:  make(chan struct{}),
		closeC: make(chan struct{}),
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
			kafka.Consumer.Unsubscribe()
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
	kafka  *internal.KafkaConsumer
	rbac   *memcached.RBAC
	events *events.RBACEvents
	doneC  chan struct{}
	closeC chan struct{}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	commit := func(msg *kafka.Message) {
		if _, err := s.kafka.Consumer.CommitMessage(msg); err != nil {
			s.logger.Error("commit failed", zap.Error(err))
		}
	}

	go func() {
		run := true

		for run {
			select {
			case <-s.closeC:
				run = false
				break
			default:
				msg, ok := s.kafka.Consumer.Poll(150).(*kafka.Message)
				if !ok {
					continue
				}

				var evt struct {
					Type  string
					Value interface{}
				}

				if err := json.NewDecoder(bytes.NewReader(msg.Value)).Decode(&evt); err != nil {
					s.logger.Info("Ignoring message, invalid", zap.Error(err))
					commit(msg)
					continue
				}

				ok = false

				switch evt.Type {
				case internaldomain.EVENT_ACCOUNT_CREATED:
					var res = evt.Value.(internaldomain.Account)
					if err := s.events.AccountCreated(res); err != nil {
						s.logger.Info("Couldn't index account", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_PROFILE_UPDATED:
					var profile = evt.Value.(internaldomain.Profile)
					if err := s.events.AccountUpdated(profile); err != nil {
						s.logger.Info("Couldn't update account", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ACCOUNT_DELETED:
					var id string = evt.Value.(string)
					if err := s.events.AccountDeleted(id); err != nil {
						s.logger.Info("Couldn't delete account", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ROLE_CREATED:
					var role = evt.Value.(internaldomain.Roles)
					if err := s.events.RoleCreated(role); err != nil {
						s.logger.Info("Couldn't index role", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ROLE_UPDATED:
					var role = evt.Value.(internaldomain.Roles)
					if err := s.events.RoleUpdated(role); err != nil {
						s.logger.Info("Couldn't update role", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ROLE_DELETED:
					var id = evt.Value.(string)
					if err := s.events.RoleDeleted(id); err != nil {
						s.logger.Info("Couldn't delete role", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_TASK_CREATED:
					var task = evt.Value.(internaldomain.Tasks)
					if err := s.events.TaskCreated(task); err != nil {
						s.logger.Info("Couldn't index task", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_TASK_UPDATED:
					var task = evt.Value.(internaldomain.Tasks)
					if err := s.events.TaskUpdated(task); err != nil {
						s.logger.Info("Couldn't update task", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_TASK_DELETED:
					var id = evt.Value.(string)
					if err := s.events.TaskDeleted(id); err != nil {
						s.logger.Info("Couldn't delete task", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ACCOUNTROLE_CREATED:
					var accountRole = evt.Value.(internaldomain.AccountRoles)
					if err := s.events.AccountRoleCreated(accountRole); err != nil {
						s.logger.Info("Couldn't index accountrole", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ACCOUNTROLE_UPDATED:
					var accountRole = evt.Value.(internaldomain.AccountRoles)
					if err := s.events.AccountRoleUpdated(accountRole); err != nil {
						s.logger.Info("Couldn't update accountrole", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ACCOUNTROLE_DELETED:
					var id = evt.Value.(string)
					if err := s.events.AccountRoleDeleted(id); err != nil {
						s.logger.Info("Couldn't delete accountrole", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ROLETASK_CREATED:
					var roleTask = evt.Value.(internaldomain.RoleTasks)
					if err := s.events.RoleTaskCreated(roleTask); err != nil {
						s.logger.Info("Couldn't index roletask", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ROLETASK_UPDATED:
					var roleTask = evt.Value.(internaldomain.RoleTasks)
					if err := s.events.RoleTaskUpdated(roleTask); err != nil {
						s.logger.Info("Couldn't update roletask", zap.Error(err))
						ok = true
					}
				case internaldomain.EVENT_ROLETASK_DELETED:
					var id = evt.Value.(string)
					if err := s.events.RoleTaskDeleted(id); err != nil {
						s.logger.Info("Couldn't delete roletask", zap.Error(err))
						ok = true
					}
				}

				if ok {
					s.logger.Info("Consumed", zap.String("type", evt.Type))
					commit(msg)
				}
			}
		}

		s.logger.Info("No more messages to consume. Exiting.")

		s.doneC <- struct{}{}
	}()

	return nil
}

// Shutdown ...
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")

	close(s.closeC)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context.Done: %w", ctx.Err())

		case <-s.doneC:
			return nil
		}
	}
}
