package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/djchanahcjd/go-rss/internal/db"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *ApiConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	// 密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error hashing password: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), db.CreateUserParams{
		ID:        uuid.New(),
		Username:  params.UserName,
		Password:  string(hashedPassword),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	respondWithJSON(w, 201, user)
}

func (apiCfg *ApiConfig) LoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		UserName string `json:"username"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}
	user, err := apiCfg.DB.GetUserByUsername(r.Context(), params.UserName)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting user: %v", err))
		return
	}
	// 密码校验
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		respondWithError(w, 400, "Incorrect password")
		return
	}
	// 登录成功，返回用户信息
	respondWithJSON(w, 200, user)
}

func (apiCfg *ApiConfig) GetUser(w http.ResponseWriter, r *http.Request, user db.User) {
	respondWithJSON(w, 200, user)
}

func (apiCfg *ApiConfig) GetPostsForUser(w http.ResponseWriter, r *http.Request, user db.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), db.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  50,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error getting posts: %v", err))
		return
	}
	respondWithJSON(w, 200, posts)
}
