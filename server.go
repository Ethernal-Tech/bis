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

	err := server.ListenAndServe()
	app.db.Close()
	log.Fatal(err)
}
