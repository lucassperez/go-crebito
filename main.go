package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lucassperez/go-crebito/applog"
	"github.com/lucassperez/go-crebito/database"
	"github.com/lucassperez/go-crebito/handlers"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s := fmt.Sprintf("\033[36m\"\033[35;1m%s\033[0;36m %s\"\033[0m", r.Method, r.URL)
			if r.Method == "POST" {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					s = fmt.Sprintf("%s %s", s, (bodyBytes))
				} else {
					s = fmt.Sprintf("%s (could not read body: %w)", s, err)
				}
				r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
			applog.WithTimeStamp(s)
			next.ServeHTTP(w, r)
		},
	)
}

func main() {
	db, closeFunc, err := database.Database()
	if err != nil {
		panic(err)
	}
	defer closeFunc()

	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("GET /clientes/{id}/extrato", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleExtrato(db, w, r)
	})
	mux.HandleFunc("POST /clientes/{id}/transacoes", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleTransacoes(db, w, r)
	})

	muxComMiddleware := logMiddleware(mux)

	porta := "4000"
	applog.WithTimeStamp("Começando o server na porta %s", porta)
	log.Fatal(http.ListenAndServe(":"+porta, muxComMiddleware))
}
