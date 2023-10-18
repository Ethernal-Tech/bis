package main

import (
	// "io/ioutil"

	"log"
	"net/http"
	"strconv"
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
	user := app.db.Login(r.Form.Get("username"), r.Form.Get("password"))

	if user == nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	app.sessionManager.Put(r.Context(), "inside", "yes")
	app.sessionManager.Put(r.Context(), "username", user.Name)
	app.sessionManager.Put(r.Context(), "bankId", user.BankId)
	app.sessionManager.Put(r.Context(), "bankName", user.BankName)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Put(r.Context(), "inside", "no")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	transactions := app.db.GetTransactionsForAddress(app.sessionManager.Get(r.Context(), "bankId").(uint64))

	viewData := map[string]any{}

	viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
	viewData["transactions"] = transactions
	viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")

	ts, err := template.ParseFiles("./static/views/home.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, viewData)
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
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	r.ParseForm()

	transactionId, _ := strconv.Atoi(r.Form.Get("transaction"))

	transaction := app.db.GetTransactionHistory(uint64(transactionId))

	viewData := map[string]any{}

	viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
	viewData["transaction"] = transaction
	viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")

	viewData["Policy301"] = "false"
	viewData["Policy707"] = "false"
	viewData["Policy17"] = "false"
	viewData["Policy444"] = "false"

	for _, policy := range transaction.Policies {
		viewData[policy] = "true"
	}

	ts, err := template.ParseFiles("./static/views/confirmtransaction.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, viewData)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) transactionHistory(w http.ResponseWriter, r *http.Request) {
	if app.sessionManager.GetString(r.Context(), "inside") != "yes" {
		http.Redirect(w, r, "/", http.StatusSeeOther)

		return
	}

	r.ParseForm()

	transactionId, _ := strconv.Atoi(r.Form.Get("transaction"))

	transaction := app.db.GetTransactionHistory(uint64(transactionId))

	viewData := map[string]any{}

	viewData["username"] = app.sessionManager.GetString(r.Context(), "username")
	viewData["transaction"] = transaction
	viewData["bankName"] = app.sessionManager.GetString(r.Context(), "bankName")

	viewData["Policy301"] = "false"
	viewData["Policy707"] = "false"
	viewData["Policy17"] = "false"
	viewData["Policy444"] = "false"

	for _, policy := range transaction.Policies {
		viewData[policy] = "true"
	}

	ts, err := template.ParseFiles("./static/views/transactionhistory.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}

	ts.Execute(w, viewData)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}
