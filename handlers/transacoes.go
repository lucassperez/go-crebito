package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func HandleTransacoes(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	id := r.PathValue("id")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{\"id\": %s, \"rota\": \"post\", \"implement\": \"me\"}\n", id)
}
