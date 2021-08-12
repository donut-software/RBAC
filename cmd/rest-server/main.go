package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rbac/cmd/internal"
	"rbac/internal/envvar"
	"rbac/internal/postgresql"
	"rbac/internal/rest"
	"rbac/internal/service"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	var env, address string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.StringVar(&address, "address", ":9234", "HTTP Server Address")
	flag.Parse()

	errC, err := run(env, address)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}
func run(env, address string) (<-chan error, error) {
	if err := envvar.Load(env); err != nil {
		return nil, fmt.Errorf("envvar.Load %w", err)
	}
	vault, err := internal.NewVaultProvider()
	if err != nil {
		return nil, fmt.Errorf("newVaultProvider %w", err)
	}
	conf := envvar.New(vault)
	db, err := internal.NewPostgreSQL(conf)
	if err != nil {
		return nil, fmt.Errorf("newDB %w", err)
	}
	srv, err := newServer(serverConfig{
		Address: address,
		Db:      db,
	})
	if err != nil {
		return nil, fmt.Errorf("newServer %w", err)
	}
	errC := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// XXX: When using Go 1.15 or older
	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc,
	// 	os.Interrupt,
	// 	syscall.SIGTERM,
	// 	syscall.SIGQUIT)

	go func() {
		// <-sc // XXX: When using Go 1.15 or older
		<-ctx.Done()

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			db.Close()
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}
	}()

	go func() {

		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
		// ErrServerClosed."
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil
}

type serverConfig struct {
	Address string
	Db      *sql.DB
}

func newServer(conf serverConfig) (*http.Server, error) {
	r := mux.NewRouter()
	repo := postgresql.NewRBAC(conf.Db)
	svc := service.NewRBAC(repo)
	rest.NewRBACHandler(svc).Register(r)
	return &http.Server{
		Handler:           r,
		Addr:              conf.Address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}, nil
}
