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

type requestParamsJSON struct {
	Valor     int    `json:"valor"`
	Tipo      string `json:"tipo"`
	Descricao string `json:"descricao"`
}

func HandleTransacoes(dbPoolChan chan *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8;")

	clienteIdStr := r.PathValue("id")
	clienteId, err := strconv.Atoi(clienteIdStr)
	if err != nil {
		unparseableId(w, clienteIdStr, err)
		return
	}

	params, ok := validateBody(w, r)
	if !ok {
		return
	}

	db := <-dbPoolChan
	defer func() {
		dbPoolChan <- db
	}()

	limite, newSaldo, err :=
		models.InsertTransacaoAndUpdateCliente(db, clienteId, params.Valor, params.Tipo, params.Descricao)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			clienteNotFound(w, clienteId, err)
			return
		} else if errors.Is(err, &models.ErrInvalidValues{}) {
			applog.WithTimeStamp("models: %s", err.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "{\"message\": \"%s\"}\n", err.Error())
			return
		} else if errors.Is(err, &models.ErrNotEnoughBalance{}) {
			applog.WithTimeStamp("models: %s", err.Error())
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, "{\"message\": \"%s\"}\n", err.Error())
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

func validateBody(w http.ResponseWriter, r *http.Request) (requestParamsJSON, bool) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		applog.WithTimeStamp("could not read body: `%s`", err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"could not read body\"\n}")
		return requestParamsJSON{}, false
	}

	var params requestParamsJSON

	err = json.Unmarshal(b, &params)
	if err != nil {
		applog.WithTimeStamp("could not unmarshall json: `%s`: %+v", b, err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"could not unmarshall json\", \"json\": %s}\n", b)
		return requestParamsJSON{}, false
	}

	json := string(b)
	if isMissingKeys(json, "valor", "tipo", "descricao") {
		applog.WithTimeStamp("json is missing keys: `%s`", string(json))
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "{\"message\": \"json missing keys\"}\n")
		return requestParamsJSON{}, false
	}

	return params, true
}
