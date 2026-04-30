package bo

import (
	"fmt"
	"time"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/enums"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/interfaces"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/models"
)

type CreateMagicLinkBO struct {
	repo         interfaces.MagicLinkRepository
	emailService interfaces.EmailService
	tokenService interfaces.TokenService
	baseURL      string
}

func NewCreateMagicLinkBO(
	repo interfaces.MagicLinkRepository,
	emailService interfaces.EmailService,
	tokenService interfaces.TokenService,
	baseURL string,
) *CreateMagicLinkBO {
	return &CreateMagicLinkBO{
		repo:         repo,
		emailService: emailService,
		tokenService: tokenService,
		baseURL:      baseURL,
	}
}

func (b *CreateMagicLinkBO) Execute(email string) error {
	token, err := b.tokenService.Generate()
	if err != nil {
		return fmt.Errorf("generate token: %w", err)
	}

	link := models.MagicLink{
		Token:     token,
		Email:     email,
		ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		Status:    enums.Pending,
	}

	if err := b.repo.Save(link); err != nil {
		return fmt.Errorf("save magic link: %w", err)
	}

	magicURL := fmt.Sprintf("%s/auth/validate?token=%s", b.baseURL, token)
	if err := b.emailService.SendMagicLink(email, magicURL); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
