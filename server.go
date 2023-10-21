package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
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

	// app.db.InsertTransactionProof(1, "Something")

	err := server.ListenAndServe()
	app.db.Close()
	log.Fatal(err)

	url := "http://localhost:9090/api/start-server"

	// Create an HTTP client.
	client := &http.Client{}

	// Create a JSON payload (in this example).
	jsonPayload := []byte(`{"tx_id": "3", "policy_id": "1"}`)

	// Create an HTTP request with the payload.
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	// Set the Content-Type header to specify the data
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Handle the response as needed.
	// (e.g., read the response body, check status codes, etc.)
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(respbody)

	//--------------------------------------------------------------------------------------
	url = "http://localhost:9090/api/start-client"

	// Create a JSON payload (in this example).
	jsonPayload = []byte(`{"tx_id": "3", "policy_id": "1", "to": "0.0.0.0:10501"}`)

	// Create an HTTP request with the payload.
	req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	// Set the Content-Type header to specify the data
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Handle the response as needed.
	// (e.g., read the response body, check status codes, etc.)
	respbody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(respbody)

}
