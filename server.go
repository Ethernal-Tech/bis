package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Failed to load .env file.")
	}

	app := &application{}

	app.dependencies()

	server := &http.Server{
		//Addr:    ":443",
		Addr:    "localhost:4000",
		Handler: app.routes(),
	}

	fmt.Println("Start server at", server.Addr)
	err = server.ListenAndServe()
	app.db.Close()
	log.Fatal(err)
}
