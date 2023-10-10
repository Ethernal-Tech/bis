package main

import (
	"log"
	"net/http"
)

func main() {
	app := &application{}

	app.dependencies()

	server := &http.Server{
		Addr:    "localhost:4000",
		Handler: app.routes(),
	}
	// trs := app.db.GetTransactionsForAddress(1)
	// fmt.Println(trs)

	// tHistory := app.db.GetTransactionHistory(1)
	// fmt.Println(tHistory)

	err := server.ListenAndServe()
	app.db.Close()
	log.Fatal(err)
}
