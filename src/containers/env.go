package main

import "os"

type config struct {
	JWTSecret   string
	BaseURL     string
	Port        string
	DynamoTable string
	SESSender   string
	SMTPHost    string
	SMTPPort    string
	SMTPFrom    string
}

func loadConfig() config {
	return config{
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
		Port:        getEnv("PORT", "8080"),
		DynamoTable: getEnv("DYNAMO_TABLE", ""),
		SESSender:   getEnv("SES_SENDER", ""),
		SMTPHost:    getEnv("SMTP_HOST", ""),
		SMTPPort:    getEnv("SMTP_PORT", "1025"),
		SMTPFrom:    getEnv("SMTP_FROM", "noreply@localhost"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
