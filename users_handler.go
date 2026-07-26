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
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		ID:          createdUser.ID,
		CreatedAt:   createdUser.CreatedAt,
		UpdatedAt:   createdUser.UpdatedAt,
		Email:       createdUser.Email,
		IsChirpyRed: createdUser.IsChirpyRed,
	})
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	request := Credentials{}
	defaultExpiry := time.Hour
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	user, err := cfg.GetUserByEmail(r.Context(), request.Email)
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
		token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, defaultExpiry)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT token", err)
			return
		}
		refreshToken := auth.MakeRefreshToken()

		err = cfg.SaveRefreshToken(r.Context(), database.SaveRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Time.AddDate(time.Now(), 0, 0, 60),
		})
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't save token", err)
			return
		}

		respondWithJSON(w, http.StatusOK, User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshToken,
			IsChirpyRed:  user.IsChirpyRed,
		})
	} else {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "An error occured while attempting to retrieve the token", nil)
		return
	}

	userID, err := auth.ValidateJwt(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to validate token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := Credentials{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user, err := cfg.UpdateEmailAndPassword(r.Context(), database.UpdateEmailAndPasswordParams{
		ID:             userID,
		Email:          request.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
