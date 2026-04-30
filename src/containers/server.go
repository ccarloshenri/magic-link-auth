package main

import (
	"log"
	"net/http"

	"magic-link-auth/src/layers/main/bo"
	"magic-link-auth/src/layers/main/controller"
	"magic-link-auth/src/layers/main/implementations/memory"
	"magic-link-auth/src/layers/main/interfaces"
	"magic-link-auth/src/layers/main/processor"
)

func main() {
	cfg := loadConfig()

	dao := memory.NewInMemoryMagicLinkDAO()
	emailService := memory.NewLogEmailService()
	tokenService := memory.NewCryptoTokenService()
	authTokenService := memory.NewJWTAuthTokenService(cfg.JWTSecret)

	createBO := bo.NewCreateMagicLinkBO(dao, emailService, tokenService, cfg.BaseURL)
	validateBO := bo.NewValidateMagicLinkBO(dao, authTokenService)

	createProcessor := processor.NewCreateMagicLinkProcessor(createBO)
	validateProcessor := processor.NewValidateMagicLinkProcessor(validateBO)

	var createCtrl interfaces.Controller = controller.NewCreateMagicLinkController(createProcessor)
	var validateCtrl interfaces.Controller = controller.NewValidateMagicLinkController(validateProcessor)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/magic-link", createCtrl.Handle)
	mux.HandleFunc("GET /auth/validate", validateCtrl.Handle)

	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}
