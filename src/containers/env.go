package main

import "os"

type config struct {
	JWTSecret string
	BaseURL   string
	Port      string
}

func loadConfig() config {
	return config{
		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		BaseURL:   getEnv("BASE_URL", "http://localhost:8080"),
		Port:      getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
