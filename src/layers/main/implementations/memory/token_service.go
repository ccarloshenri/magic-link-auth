package memory

import (
	"crypto/rand"
	"encoding/hex"
)

type CryptoTokenService struct{}

func NewCryptoTokenService() *CryptoTokenService {
	return &CryptoTokenService{}
}

func (s *CryptoTokenService) Generate() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
