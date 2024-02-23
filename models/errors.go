package models

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("resource could not be found")

type ErrNotEnoughBalance struct {
	SaldoAtual       int
	ValorDaTransacao int
	Limite           int
	ClienteID        int
}

type ErrInvalidValues struct {
	Value             string
	ValidationFailure string
}

func (e *ErrNotEnoughBalance) Error() string {
	return fmt.Sprintf(
		"cliente does not have enough balance: cliente_id: %d, limite: %d, saldo: %d, valor da transacao: %d",
		e.ClienteID, e.Limite, e.SaldoAtual, e.ValorDaTransacao,
	)
}

func (e *ErrNotEnoughBalance) Is(err error) bool {
	_, ok := err.(*ErrNotEnoughBalance)
	return ok
}

func (e *ErrInvalidValues) Error() string {
	return fmt.Sprintf("invalid params: %s %s", e.Value, e.ValidationFailure)
}

func (e *ErrInvalidValues) Is(err error) bool {
	_, ok := err.(*ErrInvalidValues)
	return ok
}
