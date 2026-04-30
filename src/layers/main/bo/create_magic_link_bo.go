package bo

import (
	"fmt"
	"time"

	"magic-link-auth/src/layers/main/enums"
	"magic-link-auth/src/layers/main/interfaces"
	"magic-link-auth/src/layers/main/models"
)

type CreateMagicLinkBO struct {
	dao          interfaces.MagicLinkDAO
	emailService interfaces.EmailService
	tokenService interfaces.TokenService
	baseURL      string
}

func NewCreateMagicLinkBO(
	dao interfaces.MagicLinkDAO,
	emailService interfaces.EmailService,
	tokenService interfaces.TokenService,
	baseURL string,
) *CreateMagicLinkBO {
	return &CreateMagicLinkBO{
		dao:          dao,
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

	if err := b.dao.Save(link); err != nil {
		return fmt.Errorf("save magic link: %w", err)
	}

	magicURL := fmt.Sprintf("%s/auth/validate?token=%s", b.baseURL, token)
	if err := b.emailService.SendMagicLink(email, magicURL); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
