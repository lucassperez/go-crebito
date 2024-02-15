package main

import (
	"fmt"
	stdlog "log"
	"net/http"
	"time"
)

func log(msg string) {
	t := time.Now()
	timeTag := t.Format("15:04:05")
	dateTag := t.Format("2006-01-02")
	fmt.Printf("[%s/%s] %s\n", timeTag, dateTag, msg)
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			s := fmt.Sprintf("\033[36m\"\033[35;1m%s\033[0;36m %s\"\033[0m", r.Method, r.URL)
			if r.Method == "POST" {
				s = fmt.Sprintf("%s %s", s, r.Body)
			}
			log(s)
			next.ServeHTTP(w, r)
		},
	)
}

func main() {
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("GET /clientes/{id}/extrato", handleExtrato)
	mux.HandleFunc("POST /clientes/{id}/transacoes", handleTransacoes)

	muxComMiddleware := logMiddleware(mux)

	porta := "4000"
	log(fmt.Sprintf("Começando o server na porta %s", porta))
	stdlog.Fatal(http.ListenAndServe(":"+porta, muxComMiddleware))
}

func handleExtrato(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	id := r.PathValue("id")
	fmt.Fprintf(w, "{\"id\": %s, \"rota\": \"get\"}", id)
}

func handleTransacoes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	id := r.PathValue("id")
	fmt.Fprintf(w, "{\"id\": %s, \"rota\": \"post\"}", id)
}
