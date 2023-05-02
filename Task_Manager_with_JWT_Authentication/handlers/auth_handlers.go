package handlers

import (
	"encoding/json"
	"net/http"

	"CRUD_REST/jwt"
	"CRUD_REST/models"
)

func Login(w http.ResponseWriter, r *http.Request) {
	// Парсим входные данные
	var userCredentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// Проверяем входные данные
	user := models.User{
		Username: "admin",
		Password: "admin",
	}

	if !user.Authenticate(userCredentials.Username, userCredentials.Password) {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Генерируем JWT-токен
	token, err := jwt.GenerateJWTToken(user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}
	respondWithJSON(w, http.StatusOK, response)
}
