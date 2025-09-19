package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/djchanahcjd/go-rss/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) CreateFeedFollows(w http.ResponseWriter, r *http.Request, user db.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	feed_follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), db.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating feed follow: %v", err))
		return
	}
	respondWithJSON(w, 201, feed_follow)
}

func (apiCfg *ApiConfig) GetFeedFollowsByUser(w http.ResponseWriter, r *http.Request, user db.User) {
	feeds, err := apiCfg.DB.GetFeedFollowsByUserID(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting feeds followed by %s: %v", user.ID, err))
		return
	}
	respondWithJSON(w, 200, feeds)
}

func (apiCfg *ApiConfig) DeleteFeedFollow(w http.ResponseWriter, r *http.Request, user db.User) {
	feedIDStr := chi.URLParam(r, "feedID")
	feedID, err := uuid.Parse(feedIDStr)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing feed_follow_id: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), db.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feedID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error deleting feed follow: %v", err))
		return
	}
	respondWithJSON(w, 200, struct{}{})
}
