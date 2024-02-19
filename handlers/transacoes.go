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

func HandleTransacoes(dbPoolChan chan *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	clienteIdStr := r.PathValue("id")
	clienteId, err := strconv.Atoi(clienteIdStr)

	if err != nil {
		unparseableId(w, clienteIdStr, err)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		applog.WithTimeStamp("could not read body: `%w`", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"could not read body\"\n}")
		return
	}

	var params requestParamsTransacaoPOST

	err = json.Unmarshal(b, &params)
	if err != nil {
		applog.WithTimeStamp("could not unmarshall json: `%s`: %+v", b, err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"could not unmarshall json\", \"json\": %s}\n", b)
		return
	}

	json := string(b)
	if isMissingKeys(json, "valor", "tipo", "descricao") {
		applog.WithTimeStamp("json is missing keys: `%s`", string(json))
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"json missing keys\"}\n")
		return
	}

	lenDescricao := len(params.Descricao)
	if params.Valor < 0 || (params.Tipo != "d" && params.Tipo != "c") || (lenDescricao < 1 || lenDescricao > 10) {
		applog.WithTimeStamp("invalid params: %+v", params)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"invalid params\"}\n")
		return
	}

	db := <-dbPoolChan
	defer func() { dbPoolChan <- db }()
	limite, newSaldo, err :=
		models.InsertTransacaoAndUpdateCliente(db, clienteId, params.Valor, params.Tipo, params.Descricao)

	if err != nil {
		if errors.Is(err, &models.ErrNotEnoughBalance{}) {
			e := err.(*models.ErrNotEnoughBalance)

			applog.WithTimeStamp(e.MoreInfo())
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "{\"message\": \"not enough balance\", \"values\": \"%s\"}\n", e.Values())
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
