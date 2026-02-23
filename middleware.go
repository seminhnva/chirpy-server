package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) showMetrics(w http.ResponseWriter, _ *http.Request) {
	res := fmt.Sprintf("Hits: %v", cfg.fileserverHits.Load())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}
