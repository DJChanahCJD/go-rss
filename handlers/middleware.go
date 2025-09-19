package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/djchanahcjd/go-rss/internal/db"
)

type authedHandler func(http.ResponseWriter, *http.Request, db.User)

// GetAPIKey extracts the API key from the request headers.
// Authorization: <api_key>
func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", fmt.Errorf("no auth header provided")
	}
	return apiKey, nil
}

func (apiCfg *ApiConfig) AuthMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err!= nil {
			log.Printf("[AUTH] Invalid API key provided: %v", err)
			respondWithError(w, 400, fmt.Sprintf("Couldn't get user: %v", err))
			return
		}
		handler(w, r, user)
	}
}

// LoggingMiddleware 记录API请求和响应的日志信息
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 获取请求信息
		method := r.Method
		path := r.URL.Path
		
		// 处理请求
		next.ServeHTTP(w, r)
		
		// 记录日志
		log.Printf(
			"[API] %s %s",
			method,
			path,
		)
	})
}