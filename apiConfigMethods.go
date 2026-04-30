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

func (cfg *apiConfig) middlewareMetricsGet() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		currentHits := cfg.fileserverHits.Load()
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		hits := fmt.Sprintf("Hits: %v", currentHits)
		w.Write([]byte(hits))
	}
}

func (cfg *apiConfig) middlewareMetricsReset(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Swap(0)
		next(w, r)
	})
}
