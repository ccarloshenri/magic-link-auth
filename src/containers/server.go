package main

import (
	"log"
	"net/http"

	"magic-link-auth/src/layers/main/bo"
	"magic-link-auth/src/layers/main/controller"
	"magic-link-auth/src/layers/main/implementations/memory"
	smtpimpl "magic-link-auth/src/layers/main/implementations/smtp"
	"magic-link-auth/src/layers/main/interfaces"
	"magic-link-auth/src/layers/main/processor"
)

func main() {
	cfg := loadConfig()

	dao := memory.NewInMemoryMagicLinkDAO()
	tokenService := memory.NewCryptoTokenService()
	authTokenService := memory.NewJWTAuthTokenService(cfg.JWTSecret)
	emailService := resolveEmailService(cfg)

	createBO := bo.NewCreateMagicLinkBO(dao, emailService, tokenService, cfg.BaseURL)
	validateBO := bo.NewValidateMagicLinkBO(dao, authTokenService)

	createProcessor := processor.NewCreateMagicLinkProcessor(createBO)
	validateProcessor := processor.NewValidateMagicLinkProcessor(validateBO)

	var createCtrl interfaces.Controller = controller.NewCreateMagicLinkController(createProcessor)
	var validateCtrl interfaces.Controller = controller.NewValidateMagicLinkController(validateProcessor)

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("static")))
	mux.HandleFunc("POST /auth/magic-link", createCtrl.Handle)
	mux.HandleFunc("GET /auth/validate", validateCtrl.Handle)

	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}

func resolveEmailService(cfg config) interfaces.EmailService {
	if cfg.SMTPHost != "" {
		log.Printf("Email: SMTP → %s:%s (from: %s)", cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom)
		return smtpimpl.NewSMTPEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom)
	}
	log.Println("Email: log-only mode (set SMTP_HOST to enable real delivery)")
	return memory.NewLogEmailService()
}
