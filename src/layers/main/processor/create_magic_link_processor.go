package processor

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/bo"
)

type CreateMagicLinkInput struct {
	Email string `json:"email"`
}

type CreateMagicLinkProcessor struct {
	bo *bo.CreateMagicLinkBO
}

func NewCreateMagicLinkProcessor(b *bo.CreateMagicLinkBO) *CreateMagicLinkProcessor {
	return &CreateMagicLinkProcessor{bo: b}
}

func (p *CreateMagicLinkProcessor) Process(input CreateMagicLinkInput) error {
	if err := validateEmail(input.Email); err != nil {
		return err
	}
	return p.bo.Execute(input.Email)
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	return nil
}
