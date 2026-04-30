package processor

import (
	"errors"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/bo"
)

type ValidateMagicLinkOutput struct {
	AccessToken string `json:"access_token"`
	Type        string `json:"type"`
}

type ValidateMagicLinkProcessor struct {
	bo *bo.ValidateMagicLinkBO
}

func NewValidateMagicLinkProcessor(b *bo.ValidateMagicLinkBO) *ValidateMagicLinkProcessor {
	return &ValidateMagicLinkProcessor{bo: b}
}

func (p *ValidateMagicLinkProcessor) Process(token string) (ValidateMagicLinkOutput, error) {
	if token == "" {
		return ValidateMagicLinkOutput{}, errors.New("token is required")
	}

	jwtToken, err := p.bo.Execute(token)
	if err != nil {
		return ValidateMagicLinkOutput{}, err
	}

	return ValidateMagicLinkOutput{
		AccessToken: jwtToken,
		Type:        "Bearer",
	}, nil
}
