package main

import (
	"log"
	"net/http"
)

type application struct {
}

func main() {
	mux := http.NewServeMux()
	app := &application{}

	mux.HandleFunc("/", app.Index)
	mux.HandleFunc("/login", app.Login)
	mux.HandleFunc("/transactions", app.Transactions)
	mux.HandleFunc("/transactions/add", app.TransactionAdd)
	mux.HandleFunc("/transactions/edit", app.TransactionEdit)
	mux.HandleFunc("/transactions/history", app.TransactionHistory)
	mux.HandleFunc("/transactions/addPolicy", app.TransactionAddPolicy)
	mux.HandleFunc("/transactions/cancel", app.TransactionCancel)
	
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
