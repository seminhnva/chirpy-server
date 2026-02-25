package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/seminhnva/chirpy-server/internal/auth"
)

func (cfg *apiConfig) handleWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil || apiKey != cfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	params := parameters{}
	if err := decodeJson(w, r, &params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode params", err)
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	rows, err := cfg.database.UpdateChirpRedById(r.Context(), params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", nil)
		return
	}
	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Couldn't find user", nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
