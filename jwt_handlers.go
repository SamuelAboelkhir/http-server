package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/SamuelAboelkhir/http-server/internal/auth"
	"github.com/SamuelAboelkhir/http-server/internal/database"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "An error occured while attempting to retrieve the refresh token", nil)
		return
	}

	token, err := cfg.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No refresh token found", nil)
		return
	}

	if token.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked", nil)
		return
	}

	if token.ExpiresAt.Compare(time.Now()) == -1 {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired", nil)
		return
	}

	newAccessToken, err := auth.MakeJWT(token.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create new access token", err)
		return
	}
	respondWithJSON(
		w, http.StatusOK, struct {
			Token string `json:"token"`
		}{
			Token: newAccessToken,
		},
	)
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "An error occured while attempting to retrieve the refresh token", nil)
		return
	}

	cfg.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		Token:     refreshToken,
	})

	w.WriteHeader(http.StatusNoContent)
}
