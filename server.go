package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	// Initialize sanction list folder
	if _, err := os.Stat("./sanction-lists"); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("sanction-lists", 0777)
			if err != nil {
				panic(fmt.Sprint("error while creating sanction-lists folder %w", err.Error()))
			}
		} else {
			panic(fmt.Sprint("error while searching for sanction-lists folder %w", err.Error()))
		}
	}

	// Initialize sanction list
	if _, err := os.Stat("./sanction-lists/UN_List.csv"); err != nil {
		if os.IsNotExist(err) {
			if _, err := getNewestSanctionsList(); err != nil {
				panic(fmt.Sprint("error while downloading latest sanction list %w", err.Error()))
			}
		} else {
			panic(fmt.Sprint("error while searching for sanction list %w", err.Error()))
		}
	}
}

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
