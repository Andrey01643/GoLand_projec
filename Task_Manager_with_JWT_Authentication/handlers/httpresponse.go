package handlers

import (
	"encoding/json"
	"log"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{Error: message}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling error response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}

// respondWithJSON отправляет HTTP-ответ в формате JSON с заданным статусом и данными.
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	jsonResponse, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshalling JSON response: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to marshal response data")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonResponse)
}
