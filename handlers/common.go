package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/djchanahcjd/go-rss/internal/db"
)

type ApiConfig struct {
	DB *db.Queries
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Server error:", msg)
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, ErrorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, 200, struct{}{})
}
