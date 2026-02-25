package main

import (
	"sync/atomic"

	"github.com/seminhnva/chirpy-server/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	database       *database.Queries
	platForm       string
	jwtSecret      string
	polka_key      string
}
