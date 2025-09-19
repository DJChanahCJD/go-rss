package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/djchanahcjd/go-rss/internal/db"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) CreateFeed(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), db.CreateFeedParams{
		ID:        uuid.New(),
		Name:      params.Name,
		Url:       params.Url,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed: %v", err))
		return
	}
	respondWithJSON(w, 201, feed)
}

func (apiCfg *ApiConfig) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetAllFeeds(r.Context())
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feeds: %v", err))
		return
	}
	respondWithJSON(w, 200, feeds)
}
