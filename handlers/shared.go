package handlers

import (
	"fmt"
	"net/http"

	"github.com/lucassperez/go-crebito/applog"
)

func clienteNotFound(w http.ResponseWriter, id int, err error) {
	applog.WithTimeStamp("cliente with id `%d` not found: %s", id, err.Error())
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "{\"message\": \"cliente not found\"}\n")
}

func somethingWentWrong(w http.ResponseWriter, err error) {
	applog.WithTimeStamp(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "{\"message\": \"something went wrong\"}\n")
}

func unparseableId(w http.ResponseWriter, id string, err error) {
	applog.WithTimeStamp("unparseable id: `%s`, %s", id, err.Error())
	w.WriteHeader(http.StatusUnprocessableEntity)
	fmt.Fprintf(w, "{\"message\": \"unparseable id\"}\n")
}
