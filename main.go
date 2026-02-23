package main

import (
	"log"
	"net/http"
)

func main() {
	apiConfig := apiConfig{}
	mux := http.NewServeMux()
	fs := (http.FileServer(http.Dir(".")))
	//readines endpoint are commonly used by external system to check out server ready to receive traffic
	mux.Handle("GET /api/healthz", http.HandlerFunc(readiness))
	mux.Handle("/app/", http.StripPrefix("/app", apiConfig.middlewareMetricsInc(fs)))

	mux.HandleFunc("GET /api/metrics", apiConfig.showMetrics)
	mux.HandleFunc("POST /api/reset", apiConfig.resetMetrics)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)
	server := &http.Server{
		Addr:    ":" + Port,
		Handler: mux,
	}
	log.Fatal(server.ListenAndServe())
}

func readiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
