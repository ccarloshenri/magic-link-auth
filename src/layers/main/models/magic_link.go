package models

import "magic-link-auth/src/layers/main/enums"

type MagicLink struct {
	Token     string
	Email     string
	ExpiresAt int64
	Status    enums.TokenStatus
}
