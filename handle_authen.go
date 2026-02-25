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

const (
	Auth   = "Authorization"
	Bearer = "Bearer"
)

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	params := parameters{}

	if err := decodeJson(w, r, &params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.database.GetUserInfoByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
	}
	isPwMatched, err := auth.CheckPasswordHash(params.Password, user.Password)
	if err != nil || !isPwMatched {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}
	refresh_token := auth.MakeRefreshToken()
	refreshToken, err := cfg.database.StoreRefreshToken(r.Context(), database.StoreRefreshTokenParams{
		Token:  refresh_token,
		UserID: user.ID,
	})

	type userResponse struct {
		Id           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}
	respondWithJSON(w, http.StatusOK, userResponse{
		Id:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken.Token,
	})
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	userInfo, err := cfg.database.GetUserFromRefreshToken(r.Context(), refresh_token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
	}
	token, err := auth.MakeJWT(userInfo.ID, cfg.jwtSecret, time.Hour)
	type res struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, http.StatusOK, res{
		Token: token,
	})
}

func (cfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	refresh_token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}

	rows, err := cfg.database.RevokeToken(r.Context(), refresh_token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", nil)
		return
	}
	if rows == 0 {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
