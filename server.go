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

	// tHistory := app.db.GetTransactionHistory(3)
	// fmt.Println(tHistory)

	// app.db.InsertTransaction(DB.Transaction{
	// 		OriginatorBank: 1,
	// 		BeneficiaryBank: 2,
	// 		Sender: 1,
	// 		Receiver: 2,
	// 		Currency: "$$",
	// 		Amount: 1234,
	// 		TypeId: 1,
	// 	})

	// user := app.db.Login("admin", "admin")
	// fmt.Println(user)

	err := server.ListenAndServe()
	app.db.Close()
	log.Fatal(err)
}
