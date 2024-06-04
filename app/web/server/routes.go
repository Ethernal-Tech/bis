package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (server *WebServer) Routes() http.Handler {
	router := httprouter.New()

	fileServerLegacy := http.FileServer(http.Dir("./static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServerLegacy))

	fileServer := http.FileServer(http.Dir("./app/web/static"))
	router.Handler(http.MethodGet, "/app/web/static/*filepath", http.StripPrefix("/app/web/static", fileServer))

	// Home controller
	{
		router.HandlerFunc(http.MethodGet, "/", server.Index)
		router.HandlerFunc(http.MethodPost, "/login", server.Login)
		router.HandlerFunc(http.MethodGet, "/logout", server.Logout)
		router.HandlerFunc(http.MethodGet, "/home", server.Home)
	}

	// Transaction controller
	{
		router.HandlerFunc(http.MethodGet, "/home/searchTransaction", server.SearchTransaction)
		router.HandlerFunc(http.MethodGet, "/addtransaction", server.AddTransaction)
		router.HandlerFunc(http.MethodPost, "/addtransaction", server.AddTransaction)
		router.HandlerFunc(http.MethodGet, "/confirmtransaction", server.ConfirmTransaction)
		router.HandlerFunc(http.MethodPost, "/confirmtransaction", server.ConfirmTransaction)
		router.HandlerFunc(http.MethodGet, "/transactionhistory", server.TransactionHistory)
	}

	// Policy controller
	{
		router.HandlerFunc(http.MethodGet, "/policies", server.ShowPolicies)
		router.HandlerFunc(http.MethodGet, "/editpolicy", server.EditPolicy)
		router.HandlerFunc(http.MethodPost, "/editpolicy", server.EditPolicy)
	}

	// Central bank controller
	{
		router.HandlerFunc(http.MethodGet, "/analytics", server.ShowAnalytics)
	}

	// API controller
	{
		router.HandlerFunc(http.MethodPost, "/api/getpolicies", server.GetPolicies)
		router.HandlerFunc(http.MethodPost, "/api/getpolicy", server.GetPolicy)
		router.HandlerFunc(http.MethodPost, "/api/submitTransactionProof", server.SubmitTransactionProof)
		router.HandlerFunc(http.MethodPost, "/api/create-transaction", server.CreateTx)
	}

	return DenyAccessToHTML(server.SessionManager.LoadAndSave(router))
}
