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
	router.HandlerFunc(http.MethodGet, "/home/searchTransaction", app.searchTransaction)
	router.HandlerFunc(http.MethodGet, "/addtransaction", app.addTransaction)
	router.HandlerFunc(http.MethodPost, "/addtransaction", app.addTransaction)
	router.HandlerFunc(http.MethodPost, "/api/getpolicies", app.getPolicies)
	router.HandlerFunc(http.MethodGet, "/policies", app.showPolicies)
	router.HandlerFunc(http.MethodGet, "/confirmtransaction", app.confirmTransaction)
	router.HandlerFunc(http.MethodPost, "/confirmtransaction", app.confirmTransaction)
	router.HandlerFunc(http.MethodGet, "/transactionhistory", app.transactionHistory)
	router.HandlerFunc(http.MethodGet, "/editpolicy", app.editPolicy)
	router.HandlerFunc(http.MethodPost, "/editpolicy", app.editPolicy)
	router.HandlerFunc(http.MethodPost, "/api/getpolicy", app.getPolicy)
	router.HandlerFunc(http.MethodGet, "/analytics", app.showAnalytics)

	router.HandlerFunc(http.MethodPost, "/api/submitTransactionProof", app.submitTransactionProof)

	return DenyAccessToHTML(app.sessionManager.LoadAndSave(router))
}
