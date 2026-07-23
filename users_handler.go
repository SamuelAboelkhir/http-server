package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SamuelAboelkhir/http-server/internal/auth"
	"github.com/SamuelAboelkhir/http-server/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Token     string    `json:"token"`
}

type Credentials struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := Credentials{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}
	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	createdUser, err := cfg.CreateUser(r.Context(), database.CreateUserParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	})
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := Credentials{}
	defaultExpiry := time.Hour
	actualExpiry := defaultExpiry
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	userExpiry := time.Duration(request.ExpiresInSeconds) * time.Second

	if userExpiry < defaultExpiry && userExpiry != 0 {
		actualExpiry = userExpiry
	}

	user, err := cfg.GetUser(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	hashMatch, err := auth.CheckPasswordHash(request.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if hashMatch {
		token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, actualExpiry)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
			return
		}
		respondWithJSON(w, http.StatusOK, User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
}
