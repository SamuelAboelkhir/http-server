package auth

import (
	"errors"
	"strings"
)

func AuthHeaderSanitizer(authHeader, prefix string) (string, error) {
	sanitized, state := strings.CutPrefix(authHeader, prefix)
	if !state {
		return "", errors.New("header malformed")
	}

	trimed := strings.TrimSpace(sanitized)

	return trimed, nil
}
