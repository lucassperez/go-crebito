package handlers

import (
	"fmt"
	"net/http"

	"github.com/lucassperez/go-crebito/log"
)

func somethingWentWrong(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "{\"message\": \"something went wrong\"}\n")
	log.WithTimeStamp(err.Error())
	return
}
