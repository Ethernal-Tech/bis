package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.index)
	router.HandlerFunc(http.MethodPost, "/login", app.login)
	router.HandlerFunc(http.MethodGet, "/logout", app.logout)
	router.HandlerFunc(http.MethodGet, "/home", app.home)

	router.HandlerFunc(http.MethodGet, "/transactions", app.transactions)
	router.HandlerFunc(http.MethodPost, "/transactions/add", app.transactionAdd)
	router.HandlerFunc(http.MethodGet, "/transactions/history", app.transactionHistory)
	router.HandlerFunc(http.MethodPost, "/transactions/addPolicy", app.transactionAddPolicy)
	router.HandlerFunc(http.MethodPost, "/transactions/cancel", app.transactionCancel)

	return app.sessionManager.LoadAndSave(router)
}
