package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Database() (db *sql.DB, closeFunc func() error, err error) {
	db, err = sql.Open("pgx", buildPgString())

	if err != nil {
		return nil, nil, fmt.Errorf("error at opening db: %w", err)
	}

	return db, db.Close, nil
}

func buildPgString() string {
	host := getEnvOrDefault("PG_HOST", "localhost")

	return fmt.Sprintf(
		"host=%s "+
			"port=5432 "+
			"user=postgres "+
			"password=postgres "+
			"dbname=rinha-go-crebito "+
			"sslmode=disable",
		host,
	)
}

func getEnvOrDefault(envName, defaultValue string) string {
	v := os.Getenv(envName)
	if v == "" {
		return defaultValue
	}
	return v
}
