package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/djchanahcjd/go-rss/internal/db"
	"github.com/google/uuid"
)

func (apiCfg * apiConfig)createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserName string `json:"username"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID: uuid.New(),
		Username: params.UserName,
		Password: "",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err!= nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}
	respondWithJSON(w, 201, user)
}

func (apiCfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJSON(w, 200, user)
}

func (apiConfig *apiConfig)getPostsForUserHandler(w http.ResponseWriter, r *http.Request, user db.User) {
	posts, err := apiConfig.DB.GetPostsForUser(r.Context(), db.GetPostsForUserParams{
		UserID: user.ID,
		Limit: 10,
	})
	if err!= nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting posts: %v", err))
		return
	}
	respondWithJSON(w, 200, posts)
}