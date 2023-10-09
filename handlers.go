package main

import (
	// "io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, app.sessionManager.GetString(r.Context(), "inside"))
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if username == "admin" && password == "password" {
		app.sessionManager.Put(r.Context(), "inside", "yes")
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Put(r.Context(), "inside", "no")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if loggedIn := app.sessionManager.GetString(r.Context(), "inside"); loggedIn != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	ts, err := template.ParseFiles("./static/views/home.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, struct{}{})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) addTransaction(w http.ResponseWriter, r *http.Request) {
	if loggedIn := app.sessionManager.GetString(r.Context(), "inside"); loggedIn != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	ts, err := template.ParseFiles("./static/views/addtransaction.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, struct{}{})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) confirmTransaction(w http.ResponseWriter, r *http.Request) {
	if loggedIn := app.sessionManager.GetString(r.Context(), "inside"); loggedIn != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

	ts, err := template.ParseFiles("./static/views/confirmtransaction.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, struct{}{})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) transactions(w http.ResponseWriter, r *http.Request) {
	// read "sender" from request
	// body, _ := ioutil.ReadAll(r.Body) ....

	// app.db.getTransactionsFor()

	// return /transactions
}

func (app *application) transactionAdd(w http.ResponseWriter, r *http.Request) {
	// Read transaction from body
	// body, _ := ioutil.ReadAll(r.Body)
	// var transactionData Transaction
	// if err := json.Unmarshal(body, &transactionData); err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	// app.db.InsertTransaction(transactionData)
	// app.db.UpdateTransactionState(TransactionStates.Initiated)

	// returect to /index
}

func (app *application) transactionHistory(w http.ResponseWriter, r *http.Request) {
	// read "transaction id" from request
	// body, _ := ioutil.ReadAll(r.Body) ....

	// app.db.GetTransactionById(transactionId)
	// app.db.GetTransactionHistoryById(transactionId)

	// return /transaction_hisotry
}

func (app *application) transactionAddPolicy(w http.ResponseWriter, r *http.Request) {
	// Read policy parameter
	// body, _ := ioutil.ReadAll(r.Body)
	// var policyData Policy
	// if err := json.Unmarshal(body, &policyData); err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	// Get transaction from db
	// app.db.GetTransactionById(policyData.transactionId)

	// Add Policy to transaction (db)
	// app.db.AddTransactionPolicy(policyData.transactionId, policyData.policies)
	// app.db.UpdateTransactionState(TransactionStates.PoliciesApplied)

	// run go routine -> Call BB Exec, Call OB Exec, Update Transaction to [ProofRequested], Wait ...
	// ... Update Transaction to [ProofReceived] or [ProofInvalid], Notify front?

	// Redirect to /index
}

func (app *application) transactionCancel(w http.ResponseWriter, r *http.Request) {
	// app.db.GetTransactionById(transactionId)
	// check if not nil
	// app.db.UpdateTransactionState(TransactionStates.Canceled)

	// Redirect to /index
}
