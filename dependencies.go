package main

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

type application struct {
	sessionManager *scs.SessionManager
}

func (app *application) dependencies() {
	sessionManager := scs.New()
	sessionManager.Store = memstore.New()
	sessionManager.Lifetime = 10 * time.Minute

	app.sessionManager = sessionManager
}
