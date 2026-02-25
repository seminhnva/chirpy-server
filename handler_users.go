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
		Id          uuid.UUID `json:"id"`
		Created_at  time.Time `json:"created_at"`
		Updated_at  time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	respondWithJSON(w, http.StatusCreated, userCreated{
		Id:          user.ID,
		Created_at:  user.CreatedAt,
		Updated_at:  user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
func (cfg *apiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", nil)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", nil)
		return
	}

	params := parameters{}
	if err := decodeJson(w, r, &params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couln't decode params", err)
		return
	}
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	rows, err := cfg.database.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:    params.Email,
		Password: hashedPassword,
		ID:       userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "something went wrong", nil)
		return
	}
	if rows == 0 {
		respondWithError(w, http.StatusUnauthorized, "unauthorized", nil)
		return
	}
	type userUpdate struct {
		Email string `json:"email"`
	}

	respondWithJSON(w, http.StatusOK, userUpdate{
		Email: params.Email,
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
