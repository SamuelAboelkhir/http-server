package main

import "net/http"

func startServer(server *http.Server) {
	server.ListenAndServe()
}
