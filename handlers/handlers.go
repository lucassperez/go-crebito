package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lucassperez/go-crebito/log"
	"github.com/lucassperez/go-crebito/models"
)

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

type extratoJSON struct {
	Saldo             saldoJSON       `json:"saldo"`
	UltimasTransacoes []transacaoJSON `json:"ultimas_transacoes"`
}

func HandleExtrato(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.WithTimeStamp(fmt.Sprintf("unparseable id: `%s`", idStr))
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"unparseable id\"\n}")
		return
	}

	cliente, timeStamp, err := models.GetCliente(db, id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			log.WithTimeStamp(fmt.Sprintf("cliente with id `%d` not found", id))
			fmt.Fprintf(w, "{\"message\": \"cliente not found\"}\n")
			return
		} else {
			somethingWentWrong(w, err)
			return
		}
	}

	transacoes, err := models.GetUltimas10Transacoes(db, cliente.ID)
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

func HandleTransacoes(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	id := r.PathValue("id")
	fmt.Fprintf(w, "{\"id\": %s, \"rota\": \"post\", \"implement\": \"me\"}\n", id)
}

func somethingWentWrong(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "{\"message\": \"something went wrong\"}\n")
	log.WithTimeStamp(err.Error())
	return
}
