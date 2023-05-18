package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"httpServer/handlers"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/test", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/test/{id}", handlers.TestHandler).Methods("GET")
	r.HandleFunc("/test/{variantID}/{taskID}/submit", handlers.NextTestHandler).Methods("POST")
	r.HandleFunc("/result", handlers.ResultTestsHandler).Methods("GET")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("POST")

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	// Обработчик для главной страницы
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/index.html")
	})
	// Запускаем сервер на порту 8080
	log.Fatal(http.ListenAndServe(":8080", r))
}
