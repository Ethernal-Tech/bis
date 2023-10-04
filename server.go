package main

import (
	"bisgo/handlers"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.Index)
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
