package interfaces

import "magic-link-auth/src/layers/main/models"

type MagicLinkDAO interface {
	Save(link models.MagicLink) error
	FindByToken(token string) (*models.MagicLink, error)
	MarkAsUsed(token string) error
}
