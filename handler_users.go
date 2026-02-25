package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/seminhnva/chirpy-server/internal/auth"
	"github.com/seminhnva/chirpy-server/internal/database"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	if err := decodeJson(w, r, &params); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	user, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		Email:    params.Email,
		Password: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating users :%v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	type userCreated struct {
		Id         uuid.UUID `json:"id"`
		Created_at time.Time `json:"created_at"`
		Updated_at time.Time `json:"updated_at"`
		Email      string    `json:"email"`
	}

	respondWithJSON(w, http.StatusCreated, userCreated{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
	})
}

func (cfg *apiConfig) handleResetUser(w http.ResponseWriter, r *http.Request) {
	if cfg.platForm != "dev" {
		respondWithError(w, http.StatusUnauthorized, "", nil)

	}
	err := cfg.database.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("Error Delete users: %v ", err)
		respondWithError(w, http.StatusInternalServerError, "", nil)
		return
	}
	w.WriteHeader(http.StatusOK)
}
