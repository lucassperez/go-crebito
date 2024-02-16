package models

import (
	"database/sql"
	"fmt"
)

type Transacao struct {
	// Nunca usado o ID?
	// ID          int
	Valor       int
	Tipo        string
	Descricao   string
	RealizadaEm string
	ClienteID   int
}

func GetUltimas10Transacoes(db *sql.DB, id_cliente int) ([]Transacao, error) {
	rows, err := db.Query(
		`SELECT valor, tipo, descricao, realizada_em FROM transacoes WHERE cliente_id = $1 LIMIT 10;`,
		id_cliente,
	)

	if err != nil {
		return nil, fmt.Errorf("get_ultimas_10_transacoes#db.Query(): %w", err)
	}

	defer rows.Close()

	ts := make([]Transacao, 10)
	var size int

	for rows.Next() {
		t := Transacao{ClienteID: id_cliente}
		err = rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &t.RealizadaEm)
		if err != nil {
			return nil, fmt.Errorf("get_ultimas_10_transacoes#rows.Scan(): %w", err)
		}
		ts[size] = t
		size++
	}

	err = rows.Err()
	if err != nil {
		// Error during iteration?
		return nil, fmt.Errorf("get_ultimas_10_transacoes#rows.Err(): %w", err)
	}

	if size < 10 {
		ts = ts[:size]
	}

	return ts, nil
}
