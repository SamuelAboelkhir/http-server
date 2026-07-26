package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/SamuelAboelkhir/http-server/internal/auth"
	"github.com/google/uuid"
)

type request struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) polkaWebhooksHandler(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to get API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Incorrect API key", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	req := request{}
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	if req.Event != "user.upgraded" {
		respondWithError(w, http.StatusNoContent, "Not an upgrade event", errors.New("not an upgrade event"))
		return
	}

	_, err = cfg.GetUserById(r.Context(), req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	err = cfg.UpgradeUserChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "User membership upgrade failed", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
