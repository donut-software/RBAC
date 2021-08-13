package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rbac/cmd/internal"
	"rbac/internal/elasticsearch"
	"rbac/internal/envvar"
	"rbac/internal/memcached"
	"rbac/internal/postgresql"
	"rbac/internal/rest"
	"rbac/internal/service"
	"syscall"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	esv7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.uber.org/zap"
)

//go:embed static
var content embed.FS

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
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("zap.NewProduction %w", err)
	}

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

	promExporter, err := internal.NewOTExporter(conf)
	if err != nil {
		return nil, fmt.Errorf("newOTExporter %w", err)
	}

	logging := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(r.Method,
				zap.Time("time", time.Now()),
				zap.String("url", r.URL.String()),
			)

			h.ServeHTTP(w, r)
		})
	}
	es, err := internal.NewElasticSearch(conf)
	if err != nil {
		return nil, fmt.Errorf("new ElasticSearch %w", err)
	}
	memcached, err := internal.NewMemcached(conf)
	if err != nil {
		return nil, fmt.Errorf("internal.NewMemcached %w", err)
	}
	srv, err := newServer(serverConfig{
		Address:       address,
		Db:            db,
		ElasticSearch: es,
		Metrics:       promExporter,
		Middlewares:   []mux.MiddlewareFunc{otelmux.Middleware("user-management-server"), logging},
		Logger:        logger,
		Memcached:     memcached,
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

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		defer func() {
			logger.Sync()
			db.Close()
			stop()
			cancel()
			close(errC)
		}()

		srv.SetKeepAlivesEnabled(false)

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving", zap.String("address", address))

		// "ListenAndServe always returns a non-nil error. After Shutdown or Close, the returned error is
		// ErrServerClosed."
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	return errC, nil
}

type serverConfig struct {
	Address       string
	Db            *sql.DB
	ElasticSearch *esv7.Client
	Metrics       http.Handler
	Middlewares   []mux.MiddlewareFunc
	Logger        *zap.Logger
	Memcached     *memcache.Client
}

func newServer(conf serverConfig) (*http.Server, error) {
	r := mux.NewRouter()
	for _, mw := range conf.Middlewares {
		r.Use(mw)
	}

	repo := postgresql.NewRBAC(conf.Db)
	search := elasticsearch.NewRBAC(conf.ElasticSearch)
	mclient := memcached.NewRBAC(conf.Memcached, search, conf.Logger)
	svc := service.NewRBAC(repo, mclient)

	rest.RegisterOpenAPI(r)
	rest.NewRBACHandler(svc).Register(r)

	fsys, _ := fs.Sub(content, "static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(fsys))))

	return &http.Server{
		Handler:           r,
		Addr:              conf.Address,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       1 * time.Second,
	}, nil
}
