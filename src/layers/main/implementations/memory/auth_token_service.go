package memory

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthTokenService struct {
	secretKey []byte
}

func NewJWTAuthTokenService(secretKey string) *JWTAuthTokenService {
	return &JWTAuthTokenService{secretKey: []byte(secretKey)}
}

func (s *JWTAuthTokenService) GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"sub": email,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}
