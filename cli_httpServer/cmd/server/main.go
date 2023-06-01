package main

import (
	"log"
	"net/http"

	"cli_httperver/server"
)

func main() {
	http.HandleFunc("/api/substring", server.SubstringHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
