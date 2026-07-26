package auth

import (
	"errors"
	"net/http"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("authorization header is missing")
	}

	token, err := AuthHeaderSanitizer(tokenString, "Bearer")
	if err != nil {
		return "", err
	}

	return token, nil
}
