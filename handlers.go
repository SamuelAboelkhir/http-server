package main

import (
	"encoding/json"
	"net/http"
)

func customHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func resetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type requestMsg struct {
		Body string `json:"body"`
	}
	type response struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	request := requestMsg{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}
	if len(request.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "A chirp can't be longer than 140 characters", nil)
		return
	}
	respondWithJSON(w, http.StatusOK, response{
		CleanedBody: checkForProfanity(request.Body),
	})
}
