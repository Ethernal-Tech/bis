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

	mux.HandleFunc("/", app.index)
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
