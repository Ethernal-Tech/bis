package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func (app *application) Index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/index.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}
	ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) Login(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/login.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}
	ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) Transactions(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/transactions.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}
	ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}	
}

func (app *application) TransactionAdd(w http.ResponseWriter, r *http.Request) {	
	ts, err := template.ParseFiles("./static/views/transaction_add.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}
	ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) TransactionEdit(w http.ResponseWriter, r *http.Request) {	
	ts, err := template.ParseFiles("./static/views/transaction_edit.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}
	ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) TransactionHistory(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/transaction_history.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}
	ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}

func (app *application) TransactionAddPolicy(w http.ResponseWriter, r *http.Request) {
	// Read policy parameter
	// body, _ := ioutil.ReadAll(r.Body)
	// var policyData Policy
	// if err := json.Unmarshal(body, &policyData); err != nil {
	// 	log.Println(err.Error())
	// 	return
	// }

	// Get transaction from db
	// app.db.transactions.get(policy.id)

	// Add Policy to transaction (db)
	
	// Redirect to transactions
}

func (app *application) TransactionCancel(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./static/views/transaction_history.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 1", 500)
		return
	}
	ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error 2", 500)
	}
}