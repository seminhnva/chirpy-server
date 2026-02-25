package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
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

func (cfg *apiConfig) handleCreateChips(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		Userid uuid.UUID `json:"user_id"`
	}
	params := parameters{}

	if err := decodeJson(w, r, &params); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirps, err := cfg.database.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body:   cleanedBody,
		UserID: params.Userid,
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
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
