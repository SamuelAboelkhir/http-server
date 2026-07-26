package auth

import "net/http"

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")

	key, err := AuthHeaderSanitizer(apiKey, "ApiKey")
	if err != nil {
		return "", err
	}

	return key, nil
}
