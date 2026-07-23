package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("authorization header is missing")
	}

	sanitized, state := strings.CutPrefix(tokenString, "Bearer")
	if !state {
		return "", errors.New("header malformed")
	}

	trimed := strings.TrimSpace(sanitized)

	return trimed, nil
}
