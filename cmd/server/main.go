package main

import (
	"log"
	"net/http"

	"github.com/AlexDillz/distributed-calculator/internal/server"
	"github.com/AlexDillz/distributed-calculator/internal/storage"
)

func main() {
	store, err := storage.NewStorage("calc.db")
	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/register", server.RegisterHandler(store))
	mux.HandleFunc("POST /api/v1/login", server.LoginHandler(store))

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
