package handlers

import (
	"CRUD_REST/db"
	"encoding/json"
	"net/http"

	"CRUD_REST/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser models.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}
	if newUser.Username == "" {
		respondWithError(w, http.StatusBadRequest, "Username is required")
		return
	}
	if newUser.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Password is required")
		return
	}

	err = db.CreateUserInDB(newUser.Username, newUser.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}
	respondWithJSON(w, http.StatusCreated, newUser)
}
