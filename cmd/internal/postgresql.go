package internal

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// NewPostgreSQL instantiates the PostgreSQL database using configuration defined in environment variables.
func NewPostgreSQL() (*sql.DB, error) {

	// XXX: We will revisit this code in future episodes replacing it with another solution
	databaseHost := "127.0.0.1"  //get("DATABASE_HOST")
	databasePort := "5432"       //get("DATABASE_PORT")
	databaseUsername := "user"   //get("DATABASE_USERNAME")
	databasePassword := "user"   //get("DATABASE_PASSWORD")
	databaseName := "rbac"       //get("DATABASE_NAME")
	databaseSSLMode := "disable" //get("DATABASE_SSLMODE")
	// XXX: -

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(databaseUsername, databasePassword),
		Host:   fmt.Sprintf("%s:%s", databaseHost, databasePort),
		Path:   databaseName,
	}

	q := dsn.Query()
	q.Add("sslmode", databaseSSLMode)

	dsn.RawQuery = q.Encode()

	db, err := sql.Open("pgx", dsn.String())
	if err != nil {
		return nil, fmt.Errorf("sql.Open %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db.Ping %w", err)
	}

	return db, nil
}
