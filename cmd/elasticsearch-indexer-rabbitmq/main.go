package main

import (
	"bytes"
	"context"
	"encoding/gob"
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

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

const rabbitMQConsumerName = "elasticsearch-indexer"

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

	rmq, err := internal.NewRabbitMQ(conf)
	if err != nil {
		return nil, fmt.Errorf("newRabbitMQ %w", err)
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

	rbEvents := events.NewRBACEvents(mClient)
	srv := &Server{
		logger: logger,
		rmq:    rmq,
		rbac:   mClient,
		events: rbEvents,
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
			rmq.Close()
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
	rmq    *internal.RabbitMQ
	rbac   *memcached.RBAC
	queue  amqp.Queue
	events *events.RBACEvents
	done   chan struct{}
}

// ListenAndServe ...
func (s *Server) ListenAndServe() error {
	// XXX: Dead Letter Exchange will be implemented in future
	q, err := s.rmq.Channel.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("channel.QueueDeclare %w", err)
	}

	err = s.rmq.Channel.QueueBind(
		q.Name,       // queue name
		"rbac.*.*.*", // routing key
		"rbac",       // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("channel.QueueBind %w", err)
	}

	msgs, err := s.rmq.Channel.Consume(
		q.Name,               // queue
		rabbitMQConsumerName, // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	if err != nil {
		return fmt.Errorf("channel.Consume %w", err)
	}

	go func() {
		for msg := range msgs {
			s.logger.Info(fmt.Sprintf("Received message: %s", msg.RoutingKey))

			var nack bool

			// XXX: Instrumentation to be added in a future
			// XXX: We will revisit defining these topics in a better way in future episodes
			switch msg.RoutingKey {
			case internaldomain.EVENT_ACCOUNT_CREATED:
				var res internaldomain.Account
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&res); err != nil {
					nack = true
					return
				}
				if err := s.events.AccountCreated(res); err != nil {
					s.logger.Info("Couldn't index account", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_PROFILE_UPDATED:
				var profile internaldomain.Profile
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&profile); err != nil {
					nack = true
					return
				}
				if err := s.events.AccountUpdated(profile); err != nil {
					s.logger.Info("Couldn't update account", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ACCOUNT_DELETED:
				var id string
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&id); err != nil {
					nack = true
					return
				}
				if err := s.events.AccountDeleted(id); err != nil {
					s.logger.Info("Couldn't delete account", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ROLE_CREATED:
				var role internaldomain.Roles
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&role); err != nil {
					nack = true
					return
				}
				if err := s.events.RoleCreated(role); err != nil {
					s.logger.Info("Couldn't index role", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ROLE_UPDATED:
				var role internaldomain.Roles
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&role); err != nil {
					nack = true
					return
				}
				if err := s.events.RoleUpdated(role); err != nil {
					s.logger.Info("Couldn't update role", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ROLE_DELETED:
				var id string
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&id); err != nil {
					nack = true
					return
				}
				if err := s.events.RoleDeleted(id); err != nil {
					s.logger.Info("Couldn't delete role", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_TASK_CREATED:
				var task internaldomain.Tasks
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&task); err != nil {
					nack = true
					return
				}
				if err := s.events.TaskCreated(task); err != nil {
					s.logger.Info("Couldn't index task", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_TASK_UPDATED:
				var task internaldomain.Tasks
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&task); err != nil {
					nack = true
					return
				}
				if err := s.events.TaskUpdated(task); err != nil {
					s.logger.Info("Couldn't update task", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_TASK_DELETED:
				var id string
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&id); err != nil {
					nack = true
					return
				}
				if err := s.events.TaskDeleted(id); err != nil {
					s.logger.Info("Couldn't delete task", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ACCOUNTROLE_CREATED:
				var accountRole internaldomain.AccountRoles
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&accountRole); err != nil {
					nack = true
					return
				}
				if err := s.events.AccountRoleCreated(accountRole); err != nil {
					s.logger.Info("Couldn't index accountrole", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ACCOUNTROLE_UPDATED:
				var accountRole internaldomain.AccountRoles
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&accountRole); err != nil {
					nack = true
					return
				}
				if err := s.events.AccountRoleUpdated(accountRole); err != nil {
					s.logger.Info("Couldn't update accountrole", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ACCOUNTROLE_DELETED:
				var id string
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&id); err != nil {
					nack = true
					return
				}
				if err := s.events.AccountRoleDeleted(id); err != nil {
					s.logger.Info("Couldn't delete accountrole", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ROLETASK_CREATED:
				var roleTask internaldomain.RoleTasks
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&roleTask); err != nil {
					nack = true
					return
				}
				if err := s.events.RoleTaskCreated(roleTask); err != nil {
					s.logger.Info("Couldn't index roletask", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ROLETASK_UPDATED:
				var roleTask internaldomain.RoleTasks
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&roleTask); err != nil {
					nack = true
					return
				}
				if err := s.events.RoleTaskUpdated(roleTask); err != nil {
					s.logger.Info("Couldn't update roletask", zap.Error(err))
					nack = true
				}
			case internaldomain.EVENT_ROLETASK_DELETED:
				var id string
				if err := gob.NewDecoder(bytes.NewReader(msg.Body)).Decode(&id); err != nil {
					nack = true
					return
				}
				if err := s.events.RoleTaskDeleted(id); err != nil {
					s.logger.Info("Couldn't delete roletask", zap.Error(err))
					nack = true
				}
			default:
				nack = true
			}

			if nack {
				s.logger.Info("NAcking :(")
				err = msg.Nack(false, nack)
			} else {
				s.logger.Info("Acking :)")
				_ = msg.Ack(false)
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

	s.rmq.Channel.Cancel(rabbitMQConsumerName, false)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context.Done: %w", ctx.Err())

		case <-s.done:
			return nil
		}
	}
}
