package auth

import (
	"errors"
	"net/http"
)

// GetAPIKey extracts the API key from the request headers.
// Authorization: <api_key>
func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("no auth header provided")
	}
	return apiKey, nil
}