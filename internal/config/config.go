package config

import (
	"os"
)

type Config struct {
	GRPCPort   string
	HTTPPort   string
	DatabaseDS string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		GRPCPort:   getEnv("GRPC_PORT", ":50051"),
		HTTPPort:   getEnv("HTTP_PORT", ":8080"),
		DatabaseDS: getEnv("DB_PATH", "calc.db"),
		JWTSecret:  getEnv("JWT_SECRET", "secret_key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
