package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/SamuelAboelkhir/http-server/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()
	dir := http.Dir(".")
	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		platform:       platform,
		Queries:        dbQueries,
	}
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(dir))))

	mux.HandleFunc("GET /api/healthz", customHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.middlewareMetricsGet())
	mux.HandleFunc("POST /admin/reset", cfg.middlewareMetricsReset(resetHandler))
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", cfg.createChirps)
	mux.HandleFunc("GET /api/chirps", cfg.fetchChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.fetchSingleChirps)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	startServer(httpServer)
}
