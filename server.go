package main

import (
	"bisgo/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/login", handlers.Login)
	mux.HandleFunc("/transactions", handlers.Transactions)
	mux.HandleFunc("/transactions/add", handlers.AddEditTransaction)
	mux.HandleFunc("/transactions/edit", handlers.AddEditTransaction)
	mux.HandleFunc("/transactions/history", handlers.TransactionHistory)

	mux.HandleFunc("/api/transactions/add", handlers.ApiTransactionAdd)
	mux.HandleFunc("/api/transactions/edit", handlers.ApiTransactionEdit)
	mux.HandleFunc("/api/transactions/addPolicy", handlers.ApiTransactionAddPolicy)
	mux.HandleFunc("/api/transactions/cancel", handlers.ApiTransactionCancel)

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
