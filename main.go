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
	jwtSecret := os.Getenv("SECRET")

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
		jwtSecret:      jwtSecret,
	}
	mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(dir))))

	mux.HandleFunc("GET /api/healthz", customHandler)
	mux.HandleFunc("GET /admin/metrics", cfg.middlewareMetricsGet())
	mux.HandleFunc("POST /admin/reset", cfg.middlewareMetricsReset(resetHandler))
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("POST /api/login", cfg.loginHandler)
	mux.HandleFunc("POST /api/chirps", cfg.createChirpsHandler)
	mux.HandleFunc("GET /api/chirps", cfg.fetchChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.fetchSingleChirpsHandler)
	mux.HandleFunc("POST /api/refresh", cfg.refresh)
	mux.HandleFunc("POST /api/revoke", cfg.revoke)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	startServer(httpServer)
}
