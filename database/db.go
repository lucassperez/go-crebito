package database

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const pgString = "host=localhost port=5432 user=postgres password=postgres dbname=rinha-go-crebito sslmode=disable"

func Database() (db *sql.DB, closeFunc func() error, err error) {
	db, err = sql.Open("pgx", pgString)

	if err != nil {
		return nil, nil, fmt.Errorf("error at opening db: %w", err)
	}
	return db, db.Close, nil
}
