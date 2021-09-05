package main

import (
	"collatz/api"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/collatz", api.CollatzHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}
