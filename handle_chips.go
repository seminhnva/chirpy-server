package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/seminhnva/chirpy-server/internal/auth"
	"github.com/seminhnva/chirpy-server/internal/database"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	type chirpResponse struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		Userid     uuid.UUID `json:"user_id"`
	}
	chirps, err := cfg.database.GetChirp(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	var chirpRes []chirpResponse
	if chirps == nil {
		chirpRes = []chirpResponse{}
	} else {
		for _, v := range chirps {
			chirpRes = append(chirpRes, chirpResponse{
				Id:         v.ID,
				Created_at: v.CreatedAt,
				Updated_at: v.UpdatedAt,
				Body:       v.Body,
				Userid:     v.UserID,
			})
		}
	}

	respondWithJSON(w, http.StatusOK, chirpRes)
}
func (cfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error(), err)
		return
	}
	chirps, err := cfg.database.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error(), err)
	}

	type chirpRes struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		Userid     uuid.UUID `json:"user_id"`
	}

	respondWithJSON(w, http.StatusOK, chirpRes{
		Id:         chirps.ID,
		Created_at: chirps.CreatedAt,
		Updated_at: chirps.UpdatedAt,
		Body:       chirps.Body,
		Userid:     chirps.UserID,
	})
}
func (cfg *apiConfig) handleDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error(), err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	chirp, err := cfg.database.GetChirpById(r.Context(), chirpID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "chirp not found", nil)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "something went wrong", nil)
		return
	}
	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You do not have permission to perform this action", nil)
		return
	}

	rows, err := cfg.database.DeleteChirpById(r.Context(), database.DeleteChirpByIdParams{
		UserID: userID,
		ID:     chirpID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", nil)
		return
	}
	if rows == 0 {
		respondWithError(w, http.StatusNotFound, "Couldn't find chirp", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handleCreateChips(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	if err := decodeJson(w, r, &params); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	chirps, err := cfg.database.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	type chirpRes struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Body       string    `json:"body"`
		Userid     uuid.UUID `json:"user_id"`
	}
	respondWithJSON(w, http.StatusCreated, chirpRes{
		Id:         chirps.ID,
		Created_at: chirps.CreatedAt,
		Updated_at: chirps.UpdatedAt,
		Body:       chirps.Body,
		Userid:     chirps.UserID,
	})
}
