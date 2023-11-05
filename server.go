package main

import (
	"fmt"
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

	fmt.Println("Start server at", server.Addr)
	err := server.ListenAndServe()
	app.db.Close()
	log.Fatal(err)
}
