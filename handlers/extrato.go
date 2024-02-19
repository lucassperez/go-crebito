package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lucassperez/go-crebito/models"
)

type extratoJSON struct {
	Saldo             saldoJSON       `json:"saldo"`
	UltimasTransacoes []transacaoJSON `json:"ultimas_transacoes"`
}

type saldoJSON struct {
	Limite      int    `json:"limite"`
	Total       int    `json:"total"`
	DataExtrato string `json:"data_extrato"`
}

type transacaoJSON struct {
	Valor       int    `json:"valor"`
	Tipo        string `json:"tipo"`
	Descricao   string `json:"descricao"`
	RealizadaEm string `json:"realizada_em"`
}

func HandleExtrato(dbPoolChan chan *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		unparseableId(w, idStr, err)
		return
	}

	db := <-dbPoolChan
	defer func() { dbPoolChan <- db }()
	cliente, timeStamp, err := models.GetCliente(db, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			clienteNotFound(w, id, err)
			return
		} else {
			somethingWentWrong(w, err)
			return
		}
	}

	transacoes, err := models.GetLast10Transacoes(db, cliente.ID)
	if err != nil {
		somethingWentWrong(w, err)
		return
	}

	var extrato extratoJSON

	extrato.Saldo.Limite = cliente.Limite
	extrato.Saldo.Total = cliente.Saldo
	extrato.Saldo.DataExtrato = timeStamp

	extrato.UltimasTransacoes = make([]transacaoJSON, len(transacoes))

	for i, t := range transacoes {
		extrato.UltimasTransacoes[i].Valor = t.Valor
		extrato.UltimasTransacoes[i].Tipo = t.Tipo
		extrato.UltimasTransacoes[i].Descricao = t.Descricao
		extrato.UltimasTransacoes[i].RealizadaEm = t.RealizadaEm
	}

	b, err := json.Marshal(extrato)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		somethingWentWrong(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, string(b))
}
