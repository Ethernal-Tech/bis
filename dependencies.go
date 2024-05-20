package main

import (
	"bisgo/db"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

type application struct {
	sessionManager *scs.SessionManager
	db             *db.DBWrapper
}

func (app *application) dependencies() {
	app.db = db.InitDb()

	sessionManager := scs.New()
	sessionManager.Store = memstore.New()
	sessionManager.Lifetime = time.Hour

	app.sessionManager = sessionManager
}
