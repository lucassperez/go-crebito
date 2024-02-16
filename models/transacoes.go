package models

import (
	"database/sql"
	"errors"
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
		`SELECT valor, tipo, descricao, realizada_em FROM transacoes `+
			`WHERE cliente_id = $1 `+
			`ORDER BY realizada_em DESC `+
			`LIMIT 10;`,
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

// InsertTransacaoAndUpdateCliente returns limite, newSaldo and error
func InsertTransacaoAndUpdateCliente(db *sql.DB, clienteId, valor int, tipo, descricao string) (int, int, error) {
	row := db.QueryRow(`SELECT limite, saldo FROM clientes WHERE id = $1;`, clienteId)

	var limite int
	var saldo int

	err := row.Scan(&limite, &saldo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, ErrNotFound
		}
		return 0, 0, fmt.Errorf("insert_transacao#row.Scan(): %w", err)
	}

	var newSaldo int

	if tipo == "d" {
		newSaldo = saldo - valor
		if newSaldo < -limite {
			return 0, saldo, &ErrNotEnoughBalance{
				SaldoAtual: saldo, ValorDaTransacao: valor, Limite: limite, ClienteID: clienteId,
			}
		}
	} else if tipo == "c" {
		newSaldo = saldo + valor
	} else {
		// TODO fazer descrição ser um tipo, um enum, algo assim
		return 0, 0, errors.New("invalid tipo")
	}

	_, err = db.Exec(
		`INSERT INTO transacoes (valor, tipo, descricao, cliente_id) VALUES ($1, $2, $3, $4)`,
		valor, tipo, descricao, clienteId,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("insert_transacao#db.Exec(transacoes): %w", err)
	}

	// Should I use returning here to return real value from the database
	// or is returning the values in these golang variables enough?
	_, err = db.Exec(
		`UPDATE clientes SET saldo = $1 WHERE id = $2;`,
		newSaldo, clienteId,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("insert_transacao#db.Exec(transacoes): %w", err)
	}

	return limite, newSaldo, nil
}
