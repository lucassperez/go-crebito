package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/lucassperez/go-crebito/applog"
	"github.com/lucassperez/go-crebito/models"
)

type requestParamsTransacaoPOST struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

func HandleTransacoes(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	clienteIdStr := r.PathValue("id")
	clienteId, err := strconv.Atoi(clienteIdStr)

	if err != nil {
		unparseableId(w, clienteIdStr, err)
		fmt.Println(clienteId)
		return
	}

	var params requestParamsTransacaoPOST

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"could not read body\"\n}")
		applog.WithTimeStamp("could not read body: `%w`", err)
		return
	}

	err = json.Unmarshal(b, &params)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"could not unmarshall json\", \"json\": %s}\n", b)
		applog.WithTimeStamp("could not unmarshall json: `%s`: %+v", b, err)
		return
	}

	json := string(b)
	if isMissingKeys(json, "valor", "tipo", "descricao") {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "{\"message\": \"json missing keys\"}\n")
		applog.WithTimeStamp("json is missing keys: `%s`", string(json))
	}

	limite, newSaldo, err :=
		models.InsertTransacaoAndUpdateCliente(db, clienteId, params.Valor, params.Tipo, params.Descricao)

	if err != nil {
		if errors.Is(err, &models.ErrNotEnoughBalance{}) {
			e := err.(*models.ErrNotEnoughBalance)

			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "{\"message\": \"not enough balance\", \"values\": \"%s\"}\n", e.Values())
			applog.WithTimeStamp(e.MoreInfo())
			return
		} else if errors.Is(err, models.ErrNotFound) {
			clienteNotFound(w, clienteId, err)
			return
		}
		somethingWentWrong(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"limite\": %d, \"saldo\": %d}\n", limite, newSaldo)
}

func isMissingKeys(json string, keys ...string) bool {
	for _, k := range keys {
		if !strings.Contains(json, fmt.Sprintf("\"%s\"", k)) {
			return true
		}
	}
	return false
}
