package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/lucassperez/go-crebito/applog"
	"github.com/lucassperez/go-crebito/database"
	"github.com/lucassperez/go-crebito/handlers"
)

func main() {
	dbPoolChan, err := database.NewDatabasePool()
	if err != nil {
		panic(err)
	}

	defer func() {
		for d := range dbPoolChan {
			d.Close()
		}
	}()

	applog.WithTimeStamp("Size of connection pool: %d", len(dbPoolChan))

	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("GET /clientes/{id}/extrato", func(w http.ResponseWriter, r *http.Request) {
		db := <-dbPoolChan
		defer func() { dbPoolChan <- db }()
		handlers.HandleExtrato(db, w, r)
	})
	mux.HandleFunc("POST /clientes/{id}/transacoes", func(w http.ResponseWriter, r *http.Request) {
		db := <-dbPoolChan
		defer func() { dbPoolChan <- db }()
		handlers.HandleTransacoes(db, w, r)
	})

	muxComMiddleware := logMiddleware(mux)

	port := os.Getenv("SERVER_ADDRESS")
	if port == "" {
		port = "4000"
		applog.WithTimeStamp("Variable SERVER_ADDRESS empty, using default value of %s", port)
	}
	applog.WithTimeStamp("Starting the server at port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, muxComMiddleware))
}

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
