package main

import (
	"bisgo/DB"
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

	db := DB.InitDb()

	err := server.ListenAndServe()
	db.Close()
	log.Fatal(err)
}
