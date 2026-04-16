package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	dir := http.Dir(".")
	mux.Handle("/", http.FileServer(dir))

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	startServer(httpServer)
}
