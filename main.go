package main

import (
	"collatz/api"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/collatz", api.CollatzHandler)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
