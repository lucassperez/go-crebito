package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// 50 is max_connections in postgresql.conf and
// 2 are the number of instances I'll have running in prod
const poolSize int = 50 / 2

func NewDatabasePool() (dbPoolChan chan *sql.DB, err error) {
	dbPoolChan = make(chan *sql.DB, poolSize)

	for i := 1; i <= poolSize; i++ {
		db, err := sql.Open("pgx", buildPgString())
		if err != nil {
			return nil, fmt.Errorf("error at opening db #%d: %w", i, err)
		}
		dbPoolChan <- db
	}

	return dbPoolChan, nil
}

func buildPgString() string {
	host := getEnvOrDefault("DB_HOST", "localhost")
	user := getEnvOrDefault("DB_USER", "rinheiro")
	password := getEnvOrDefault("DB_PASSWORD", "rinha123")
	dbname := getEnvOrDefault("DB_NAME", "rinha-go-crebito")

	return fmt.Sprintf(
		"host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
		host, user, password, dbname,
	)
}

func getEnvOrDefault(envName, defaultValue string) string {
	v := os.Getenv(envName)
	if v == "" {
		return defaultValue
	}
	return v
}
