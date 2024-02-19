package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Cliente struct {
	ID     int
	Limite int
	Saldo  int
}

func GetCliente(db *sql.DB, id int) (*Cliente, string, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, "", fmt.Errorf("get_cliente: could not start transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT limite, saldo, LOCALTIMESTAMP FROM clientes WHERE id = $1 FOR NO KEY UPDATE;`, id)

	c := Cliente{ID: id}
	var timeStamp string

	err = row.Scan(&c.Limite, &c.Saldo, &timeStamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("get_cliente: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, "", fmt.Errorf("get_cliente: could not commit the transaction: %w", err)
	}

	return &c, timeStamp, nil
}
