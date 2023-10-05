package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

type application struct {
	sessionManager *scs.SessionManager
}

func main() {
	sessionManager := scs.New()
	sessionManager.Store = memstore.New()
	sessionManager.Lifetime = 2 * time.Hour

	app := &application{
		sessionManager: sessionManager,
	}

	server := &http.Server{
		Addr:    "localhost:4000",
		Handler: app.routes(),
	}

	err := server.ListenAndServe()
	log.Fatal(err)
}
