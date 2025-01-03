package main

import (
	"log"
	"net/http"
)

func main() {
	h := NewHandler()
	go h.manager.Run()

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	mux.HandleFunc("POST /create", h.CreateGroup)
	mux.HandleFunc("DELETE /groups/{id}", h.DeleteGroup)
	mux.HandleFunc("GET /groups/{id}", h.ServeWS)

	log.Fatal(server.ListenAndServe())
}
