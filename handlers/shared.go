package handlers

import (
	"fmt"
	"net/http"

	"github.com/lucassperez/go-crebito/applog"
)

func somethingWentWrong(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "{\"message\": \"something went wrong\"}\n")
	applog.WithTimeStamp(err.Error())
	return
}

func unparseableId(w http.ResponseWriter, id string, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	fmt.Fprintf(w, "{\"message\": \"unparseable id\"\n}")
	applog.WithTimeStamp("unparseable id: `%s`", id)
}
