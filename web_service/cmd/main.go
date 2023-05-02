package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"web_service/handlers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/get-items/{id}", handlers.GetItemsHandler)
	http.ListenAndServe(":8080", r)
}
