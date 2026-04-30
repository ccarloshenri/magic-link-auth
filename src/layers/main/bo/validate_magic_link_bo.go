package bo

import (
	"errors"
	"fmt"
	"time"

	"github.com/carlos-sousa/magic-link-auth/src/layers/main/enums"
	"github.com/carlos-sousa/magic-link-auth/src/layers/main/interfaces"
)

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")
	ErrTokenUsed     = errors.New("token already used")
)

type ValidateMagicLinkBO struct {
	repo             interfaces.MagicLinkRepository
	authTokenService interfaces.AuthTokenService
}

func NewValidateMagicLinkBO(
	repo interfaces.MagicLinkRepository,
	authTokenService interfaces.AuthTokenService,
) *ValidateMagicLinkBO {
	return &ValidateMagicLinkBO{repo: repo, authTokenService: authTokenService}
}

func (b *ValidateMagicLinkBO) Execute(token string) (string, error) {
	link, err := b.repo.FindByToken(token)
	if err != nil {
		return "", ErrTokenNotFound
	}

	if link.Status == enums.Used {
		return "", ErrTokenUsed
	}

	if link.Status == enums.Expired || time.Now().Unix() > link.ExpiresAt {
		return "", ErrTokenExpired
	}

	if err := b.repo.MarkAsUsed(token); err != nil {
		return "", fmt.Errorf("mark token as used: %w", err)
	}

	jwtToken, err := b.authTokenService.GenerateJWT(link.Email)
	if err != nil {
		return "", fmt.Errorf("generate JWT: %w", err)
	}

	return jwtToken, nil
}
