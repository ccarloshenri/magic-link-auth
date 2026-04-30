package main

import (
	"context"
	"log"
	"net/http"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"

	"magic-link-auth/src/layers/main/bo"
	"magic-link-auth/src/layers/main/controller"
	awsimpl "magic-link-auth/src/layers/main/implementations/aws"
	"magic-link-auth/src/layers/main/implementations/memory"
	smtpimpl "magic-link-auth/src/layers/main/implementations/smtp"
	"magic-link-auth/src/layers/main/interfaces"
	"magic-link-auth/src/layers/main/processor"
)

func main() {
	cfg := loadConfig()

	dao := resolveDAO(cfg)
	emailService := resolveEmailService(cfg)
	tokenService := memory.NewCryptoTokenService()
	authTokenService := memory.NewJWTAuthTokenService(cfg.JWTSecret)

	createBO := bo.NewCreateMagicLinkBO(dao, emailService, tokenService, cfg.BaseURL)
	validateBO := bo.NewValidateMagicLinkBO(dao, authTokenService)

	createProcessor := processor.NewCreateMagicLinkProcessor(createBO)
	validateProcessor := processor.NewValidateMagicLinkProcessor(validateBO)

	var createCtrl interfaces.Controller = controller.NewCreateMagicLinkController(createProcessor)
	var validateCtrl interfaces.Controller = controller.NewValidateMagicLinkController(validateProcessor)

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServer(http.Dir("ui")))
	mux.HandleFunc("POST /auth/magic-link", createCtrl.Handle)
	mux.HandleFunc("GET /auth/validate", validateCtrl.Handle)

	log.Printf("Server running on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}

func resolveDAO(cfg config) interfaces.MagicLinkDAO {
	if cfg.DynamoTable != "" {
		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Fatalf("load AWS config: %v", err)
		}
		log.Printf("DAO: DynamoDB → table %s", cfg.DynamoTable)
		return awsimpl.NewDynamoDBMagicLinkDAO(dynamodb.NewFromConfig(awsCfg), cfg.DynamoTable)
	}
	log.Println("DAO: in-memory (set DYNAMO_TABLE to use DynamoDB)")
	return memory.NewInMemoryMagicLinkDAO()
}

func resolveEmailService(cfg config) interfaces.EmailService {
	if cfg.SMTPHost != "" {
		log.Printf("Email: SMTP → %s:%s (from: %s)", cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom)
		return smtpimpl.NewSMTPEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPFrom)
	}
	if cfg.SESSender != "" {
		awsCfg, err := awsconfig.LoadDefaultConfig(context.Background())
		if err != nil {
			log.Fatalf("load AWS config: %v", err)
		}
		log.Printf("Email: SES → from %s", cfg.SESSender)
		return awsimpl.NewSESEmailService(sesv2.NewFromConfig(awsCfg), cfg.SESSender)
	}
	log.Println("Email: log-only mode (set SMTP_HOST or SES_SENDER to enable delivery)")
	return memory.NewLogEmailService()
}
