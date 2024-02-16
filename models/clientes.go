package models

import (
	"database/sql"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("models: resource could not be found")

type Cliente struct {
	ID     int
	Limite int
	Saldo  int
}

func GetCliente(db *sql.DB, id int) (*Cliente, error) {
	row := db.QueryRow(`SELECT limite, saldo FROM clientes WHERE id = $1;`, id)

	c := Cliente{ID: id}

	err := row.Scan(&c.Limite, &c.Saldo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("get_cliente: %w", err)
	}

	return &c, nil
}
