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
		router.HandlerFunc(http.MethodGet, "/", server.Login)
		router.HandlerFunc(http.MethodPost, "/login", server.Login)
		router.HandlerFunc(http.MethodGet, "/logout", server.Logout)
		router.HandlerFunc(http.MethodGet, "/home", server.Home)
	}

	// TODO: delete everything related to transaction controller (transaction -> compilance check)

	// Transaction controller
	{
		router.HandlerFunc(http.MethodPost, "/transactions", server.GetTransactions)
		router.HandlerFunc(http.MethodGet, "/addtransaction", server.AddTransaction)
		router.HandlerFunc(http.MethodPost, "/addtransaction", server.AddTransaction)
		router.HandlerFunc(http.MethodGet, "/confirmtransaction", server.ConfirmTransaction)
		router.HandlerFunc(http.MethodPost, "/confirmtransaction", server.ConfirmTransaction)
		router.HandlerFunc(http.MethodGet, "/transactionhistory", server.TransactionHistory)
	}

	// Compliance check controller
	{
		router.HandlerFunc(http.MethodGet, "/complianceCheckIndex", server.ComplianceCheckIndex)
		router.HandlerFunc(http.MethodPost, "/compliancechecks", server.ComplianceChecks)
		router.HandlerFunc(http.MethodGet, "/addcompliancecheck", server.AddComplianceCheck)
		router.HandlerFunc(http.MethodPost, "/addcompliancecheck", server.AddComplianceCheck)
		router.HandlerFunc(http.MethodGet, "/compliancecheckdetails", server.ComplianceCheckDetails)
		router.HandlerFunc(http.MethodGet, "/confirmcompliancecheck", server.ConfirmComplianceCheck)
		router.HandlerFunc(http.MethodPost, "/confirmcompliancecheck", server.ConfirmComplianceCheck)
	}

	// Policy controller
	{
		router.HandlerFunc(http.MethodPost, "/policies", server.Policies)
		router.HandlerFunc(http.MethodGet, "/editpolicy", server.EditPolicy)
		router.HandlerFunc(http.MethodPost, "/editpolicy", server.EditPolicy)
		router.HandlerFunc(http.MethodGet, "/addpolicy", server.AddPolicyGetModel)
	}

	// Bank controller
	{
		router.HandlerFunc(http.MethodPost, "/analytics", server.Analytics)
	}

	// Central bank controller
	// {
	// 	router.HandlerFunc(http.MethodGet, "/analytics", server.ShowAnalytics)
	// }

	// API controller
	{
		router.HandlerFunc(http.MethodPost, "/api/getbeneficiarybankpolicies", server.GetBeneficiaryBankPolicies)
		router.HandlerFunc(http.MethodPost, "/api/getpolicy", server.GetPolicy)
	}

	return DenyAccessToHTML(server.SessionManager.LoadAndSave(router))
}
