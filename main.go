package main

import (
	"net/http"
	"sync/atomic"
)

func main() {
	mux := http.NewServeMux()
	dir := http.Dir(".")
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(dir))))

	mux.HandleFunc("GET /api/healthz", customHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.middlewareMetricsGet())
	mux.HandleFunc("POST /admin/reset", cfg.middlewareMetricsReset(resetHandler))
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	startServer(httpServer)
}
