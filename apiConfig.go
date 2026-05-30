package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/SamuelAboelkhir/http-server/internal/database"
	"github.com/google/uuid"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	platform       string
	*database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetricsGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		currentHits := cfg.fileserverHits.Load()
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		hits := fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, currentHits)
		w.Write([]byte(hits))
	}
}

func (cfg *apiConfig) middlewareMetricsReset(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	platform := cfg.platform
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		cfg.ResetUsers(r.Context())
		cfg.fileserverHits.Swap(0)
		next(w, r)
	})
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type requestMsg struct {
		Email string `json:"email"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	request := requestMsg{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	createdUser, err := cfg.CreateUser(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	})
}

func (cfg *apiConfig) chirps(w http.ResponseWriter, r *http.Request) {
	type requestMsg struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	request := requestMsg{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	isValid := validateChirpHandler(w, request.Body)

	if isValid {
		createdChirp, err := cfg.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   request.Body,
			UserID: request.UserID,
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create chirp", err)
		}

		respondWithJSON(w, http.StatusCreated, response{
			ID:        createdChirp.ID,
			CreatedAt: createdChirp.CreatedAt,
			UpdatedAt: createdChirp.UpdatedAt,
			Body:      createdChirp.Body,
			UserID:    createdChirp.UserID,
		})
	}
}
