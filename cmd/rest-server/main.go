package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"rbac/cmd/internal"
	"rbac/internal/postgresql"
	"rbac/internal/rest"
	"rbac/internal/service"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	db, err := internal.NewPostgreSQL()
	if err != nil {
		fmt.Println(fmt.Errorf("newDB %w", err))
	}
	srv, err := newServer(serverConfig{
		Db:      db,
		Address: ":9234",
	})
	if err != nil {
		fmt.Println(fmt.Errorf("new server %w", err))
	}
	srv.ListenAndServe()
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
