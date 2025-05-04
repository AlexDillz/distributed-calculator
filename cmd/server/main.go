package main

import (
	"log"
	"net/http"

	"github.com/AlexDillz/distributed-calculator/internal/config"
	"github.com/AlexDillz/distributed-calculator/internal/server"
	"github.com/AlexDillz/distributed-calculator/internal/storage"
	"github.com/AlexDillz/distributed-calculator/pkg/logging"
)

func main() {
	cfg := config.Load()
	logging.InitLogger()
	logger := logging.GetLogger()

	store, err := storage.NewStorage(cfg.DatabaseDS)
	if err != nil {
		logger.Fatalf("Failed to connect to DB: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/register", server.RegisterHandler(store))
	mux.HandleFunc("/api/v1/login", server.LoginHandler(store))
	mux.HandleFunc("/api/v1/calculate",
		server.AuthMiddleware(server.CalculateHandler(store)),
	)
	mux.HandleFunc("/api/v1/expressions",
		server.AuthMiddleware(server.ListExpressionsHandler(store)),
	)
	mux.HandleFunc("/api/v1/expressions/",
		server.AuthMiddleware(server.GetExpressionHandler(store)),
	)

	logger.Printf("HTTP Server running on %s\n", cfg.HTTPPort)
	log.Fatal(http.ListenAndServe(cfg.HTTPPort, mux))
}
