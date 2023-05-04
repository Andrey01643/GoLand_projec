package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"CRUD_REST/db"
	"CRUD_REST/handlers"
)

func main() {
	dataBase, err := db.OpenDB()
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}
	defer dataBase.Close()

	router := mux.NewRouter()

	router.HandleFunc("/users", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/tasks", handlers.GetTasks).Methods("GET")
	router.HandleFunc("/task/{id}", handlers.GetTask).Methods("GET")
	router.HandleFunc("/tasks", handlers.CreateTask).Methods("POST")
	router.HandleFunc("/tasks/{id}", handlers.UpdateTask).Methods("PUT")
	router.HandleFunc("/tasks/{id}", handlers.DeleteTask).Methods("DELETE")
	router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetAPIDocs(w, router)
	}).Methods("GET")

	errLis := http.ListenAndServe(":8000", router)
	if errLis != nil {
		log.Fatal("Failed to start server: ", err)
	}
}
