package main

import (
	"bisgo/DB"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

type application struct {
	sessionManager *scs.SessionManager
	db             *DB.DBWrapper
}

func (app *application) dependencies() {
	app.db = DB.InitDb()

	sessionManager := scs.New()
	sessionManager.Store = memstore.New()
	sessionManager.Lifetime = 10 * time.Minute

	app.sessionManager = sessionManager
}
