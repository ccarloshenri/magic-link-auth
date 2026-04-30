package interfaces

import "github.com/carlos-sousa/magic-link-auth/src/layers/main/models"

type MagicLinkRepository interface {
	Save(link models.MagicLink) error
	FindByToken(token string) (*models.MagicLink, error)
	MarkAsUsed(token string) error
}
