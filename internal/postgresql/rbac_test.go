package postgresql_test

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"rbac/internal"
	"rbac/internal/postgresql"
	"testing"
	"time"

	migrate "github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/go-cmp/cmp"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

func createAcc() internal.Account {
	profile := internal.Profile{
		Profile_Picture:    "image",
		Profile_Background: "image",
		First_Name:         "test",
		Last_Name:          "test",
		Mobile:             "091234567891",
		Email:              "test@test.com",
	}
	acc := internal.Account{
		Profile:  profile,
		UserName: "test",
	}
	return acc
}
func TestAccount_Create(t *testing.T) {
	t.Parallel()

	t.Run("Create: OK", func(t *testing.T) {
		t.Parallel()

		err := postgresql.NewRBAC(newDB(t)).CreateAccount(context.Background(), createAcc(), "test")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	})

	// t.Run("Create: Error Email", func(t *testing.T) {
	// 	t.Parallel()
	// 	err := postgresql.NewRBAC(newDB(t)).CreateAccount(context.Background(), createAcc(), "test")
	// 	if err != nil { // invalid email
	// 		t.Fatalf("expected error, got %s", err)
	// 	}
	// })

	t.Run("Find account: Ok", func(t *testing.T) {
		t.Parallel()
		acc := createAcc()
		db := newDB(t)
		err := postgresql.NewRBAC(db).CreateAccount(context.Background(), acc, "test")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		account, err := postgresql.NewRBAC(db).Account(context.Background(), "test")
		acc.Id = account.Id
		acc.Profile.Id = account.Profile.Id
		acc.Profile.CreatedAt = account.Profile.CreatedAt
		acc.CreatedAt = account.CreatedAt
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		if !cmp.Equal(acc, account) {
			t.Fatalf("expected result does not match: %s", cmp.Diff(acc, account))
		}
	})
	t.Run("Change Password: Ok", func(t *testing.T) {
		t.Parallel()
		acc := createAcc()
		db := newDB(t)
		err := postgresql.NewRBAC(db).CreateAccount(context.Background(), acc, "test")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
		err = postgresql.NewRBAC(db).ChangePassword(context.Background(), acc.UserName, "abc123")
		if err != nil {
			t.Fatalf("expected no error, got %s", err)
		}
	})

	// t.Run("Delete Account: Ok", func(t *testing.T) {
	// 	t.Parallel()
	// 	acc := createAcc()
	// 	db := newDB(t)
	// 	err := postgresql.NewRBAC(db).CreateAccount(context.Background(), acc, "test")
	// 	if err != nil {
	// 		t.Fatalf("expected no error, got %s", err)
	// 	}
	// 	account, err := postgresql.NewRBAC(db).Account(context.Background(), "test")
	// 	if err != nil {
	// 		t.Fatalf("expected no error, got %s", err)
	// 	}
	// 	if account.IsBlocked == false {
	// 		t.Fatalf("expected false, got %t", account.IsBlocked)
	// 	}
	// })
}

func newDB(tb testing.TB) *sql.DB {
	dsn := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword("root", "root"),
		Path:   "user-management",
	}

	q := dsn.Query()
	q.Add("sslmode", "disable")

	dsn.RawQuery = q.Encode()

	//-

	pool, err := dockertest.NewPool("")
	if err != nil {
		tb.Fatalf("Couldn't connect to docker: %s", err)
	}

	pool.MaxWait = 10 * time.Second

	pw, _ := dsn.User.Password()

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13.3-alpine",
		Env: []string{
			fmt.Sprintf("POSTGRES_USER=%s", dsn.User.Username()),
			fmt.Sprintf("POSTGRES_PASSWORD=%s", pw),
			fmt.Sprintf("POSTGRES_DB=%s", dsn.Path),
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})

	if err != nil {
		tb.Fatalf("Couldn't start resource: %s", err)
	}

	resource.Expire(60)

	tb.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			tb.Fatalf("Couldn't purge container: %v", err)
		}
	})

	dsn.Host = fmt.Sprintf("%s:5432", resource.Container.NetworkSettings.IPAddress)
	// if runtime.GOOS == "darwin" { // MacOS-specific
	// 	dsn.Host = net.JoinHostPort(resource.GetBoundIP("1111/tcp"), resource.GetPort("5432/tcp"))
	// }

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		tb.Fatalf("Couldn't open DB: %s", err)
	}

	tb.Cleanup(func() {
		if err := db.Close(); err != nil {
			tb.Fatalf("Couldn't close DB: %s", err)
		}
	})

	if err := pool.Retry(func() (err error) {
		return db.Ping()
	}); err != nil {
		tb.Fatalf("Couldn't ping DB: %s", err)
	}

	//-

	instance, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		tb.Fatalf("Couldn't migrate (1): %s", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://../../db/migrations/", "postgres", instance)
	if err != nil {
		tb.Fatalf("Couldn't migrate (2): %s", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		tb.Fatalf("Couldnt' migrate (3): %s", err)
	}

	return db
}
