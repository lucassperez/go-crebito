package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Transacao struct {
	Valor       int
	Tipo        string
	Descricao   string
	RealizadaEm string
	ClienteID   int
}

func GetLast10Transacoes(db *sql.DB, id_cliente int) ([]Transacao, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("get_last_10_transacoes: could not start transaction: %w", err)
	}
	defer tx.Rollback()

	rows, err := tx.Query(
		`SELECT valor, tipo, descricao, realizada_em FROM transacoes `+
			`WHERE cliente_id = $1 `+
			`ORDER BY realizada_em DESC `+
			`LIMIT 10 `+
			`FOR NO KEY UPDATE;`,
		id_cliente,
	)

	if err != nil {
		return nil, fmt.Errorf("get_last_10_transacoes#tx.Query(): %w", err)
	}

	defer rows.Close()

	ts := make([]Transacao, 10)
	var size int

	for rows.Next() {
		t := Transacao{ClienteID: id_cliente}
		err = rows.Scan(&t.Valor, &t.Tipo, &t.Descricao, &t.RealizadaEm)
		if err != nil {
			return nil, fmt.Errorf("get_last_10_transacoes#rows.Scan(): %w", err)
		}
		ts[size] = t
		size++
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get_last_10_transacoes#rows.Err(): %w", err)
	}

	if size < 10 {
		ts = ts[:size]
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("get_last_10_transacoes: could not commit the transaction: %w", err)
	}

	return ts, nil
}

// InsertTransacaoAndUpdateCliente returns limite, newSaldo and error
func InsertTransacaoAndUpdateCliente(db *sql.DB, clienteId, valor int, tipo, descricao string) (int, int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, 0, fmt.Errorf("insert_transacao: could not start transaction: %w", err)
	}
	defer tx.Rollback()

	row := tx.QueryRow(`SELECT limite, saldo FROM clientes WHERE id = $1 FOR NO KEY UPDATE;`, clienteId)

	var limite int
	var saldo int

	err = row.Scan(&limite, &saldo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, fmt.Errorf("insert_transacao: cliente not found: %w", ErrNotFound)
		}
		return 0, 0, fmt.Errorf("insert_transacao#row.Scan(): %w", err)
	}

	var newSaldo int

	if tipo == "d" {
		newSaldo = saldo - valor
	} else if tipo == "c" {
		newSaldo = saldo + valor
	}

	if newSaldo < (limite * -1) {
		return 0, saldo, &ErrNotEnoughBalance{
			SaldoAtual: saldo, ValorDaTransacao: valor, Limite: limite, ClienteID: clienteId,
		}
	}

	_, err = tx.Exec(
		`INSERT INTO transacoes (valor, tipo, descricao, cliente_id) VALUES ($1, $2, $3, $4);`,
		valor, tipo, descricao, clienteId,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("insert_transacao#tx.Exec(insert transacoes): %w", err)
	}

	_, err = tx.Exec(
		`UPDATE clientes SET saldo = $1 WHERE id = $2;`,
		newSaldo, clienteId,
	)
	if err != nil {
		return 0, 0, fmt.Errorf("insert_transacao#tx.Exec(update clientes transacoes): %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, 0, fmt.Errorf("insert_transacao: could not commit the transaction: %w", err)
	}

	return limite, newSaldo, nil
}
