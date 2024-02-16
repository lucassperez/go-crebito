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
	row := db.QueryRow(`SELECT limite, saldo, LOCALTIMESTAMP FROM clientes WHERE id = $1;`, id)

	c := Cliente{ID: id}
	var timeStamp string

	err := row.Scan(&c.Limite, &c.Saldo, &timeStamp)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, "", ErrNotFound
		}
		return nil, "", fmt.Errorf("get_cliente: %w", err)
	}

	return &c, timeStamp, nil
}
