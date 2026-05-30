package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	type errorResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}

func checkForProfanity(body string) string {
	profanities := map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}
	words := strings.Fields(body)
	for i, word := range words {
		normalizedWord := strings.ToLower(word)
		if _, ok := profanities[normalizedWord]; ok {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}

func validateChirpHandler(w http.ResponseWriter, payload string) bool {
	type requestMsg struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	if len(payload) > 140 {
		respondWithError(w, http.StatusBadRequest, "A chirp can't be longer than 140 characters", nil)
		return false
	}

	return true
}
