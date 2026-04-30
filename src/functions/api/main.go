package main

import (
	"log"
	"net/http"
	"os"

	"github.com/carlos-sousa/magic-link-auth/src/functions/api/handler"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/bo"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/implementations/memory"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/processor"
)

func main() {
	jwtSecret := getEnv("JWT_SECRET", "dev-secret-change-in-production")
	baseURL := getEnv("BASE_URL", "http://localhost:8080")

	repo := memory.NewInMemoryMagicLinkRepository()
	emailService := memory.NewLogEmailService()
	tokenService := memory.NewCryptoTokenService()
	authTokenService := memory.NewJWTAuthTokenService(jwtSecret)

	createBO := bo.NewCreateMagicLinkBO(repo, emailService, tokenService, baseURL)
	validateBO := bo.NewValidateMagicLinkBO(repo, authTokenService)

	createProcessor := processor.NewCreateMagicLinkProcessor(createBO)
	validateProcessor := processor.NewValidateMagicLinkProcessor(validateBO)

	createHandler := handler.NewCreateMagicLinkHandler(createProcessor)
	validateHandler := handler.NewValidateMagicLinkHandler(validateProcessor)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/magic-link", createHandler.Handle)
	mux.HandleFunc("GET /auth/validate", validateHandler.Handle)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
