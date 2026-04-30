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
	mux.Handle("/app", http.StripPrefix("/app", cfg.middlewareMetricsInc(http.FileServer(dir))))

	mux.HandleFunc("/healthz", customHandler)
	mux.HandleFunc("/metrics", cfg.middlewareMetricsGet())
	mux.HandleFunc("/reset", cfg.middlewareMetricsReset(resetHandler))

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	startServer(httpServer)
}
