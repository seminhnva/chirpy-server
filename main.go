package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/seminhnva/chirpy-server/internal/database"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platForm := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")
	if dbURL == "" {
		log.Fatal("DB_URL is empty")
	}
	if platForm == "" {
		log.Fatal("PLATFORM must be set")
	}
	if jwtSecret == "" {
		log.Fatal("PLATFORM must be set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("cannot open db: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	dbQueries := database.New(db)
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		database:       dbQueries,
		platForm:       platForm,
		jwtSecret:      jwtSecret,
	}
	mux := http.NewServeMux()
	fs := (http.FileServer(http.Dir(".")))
	//readines endpoint are commonly used by external system to check out server ready to receive traffic
	mux.Handle("/app/", http.StripPrefix("/app", apiCfg.middlewareMetricsInc(fs)))

	mux.Handle("GET /api/healthz", http.HandlerFunc(handlerReadiness))
	mux.HandleFunc("GET /api/metrics", apiCfg.showMetrics)

	mux.HandleFunc("POST /admin/reset", apiCfg.handleResetUser)
	mux.HandleFunc("POST /api/reset", apiCfg.resetMetrics)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)

	mux.HandleFunc("POST /api/chirps", apiCfg.handleCreateChips)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirpByID)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevokeToken)
	server := &http.Server{
		Addr:    ":" + Port,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}
