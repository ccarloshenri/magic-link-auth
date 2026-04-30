package testutil

import (
	"errors"

	"magic-link-auth/src/layers/main/models"
)

type MockTokenService struct {
	Token string
	Err   error
}

func (m *MockTokenService) Generate() (string, error) {
	return m.Token, m.Err
}

type MockEmailService struct {
	Err        error
	CalledWith []string
}

func (m *MockEmailService) SendMagicLink(email, link string) error {
	m.CalledWith = append(m.CalledWith, email)
	return m.Err
}

type MockMagicLinkDAO struct {
	SaveErr       error
	FindResult    *models.MagicLink
	FindErr       error
	MarkAsUsedErr error
	MarkedAsUsed  []string
}

func (m *MockMagicLinkDAO) Save(_ models.MagicLink) error {
	return m.SaveErr
}

func (m *MockMagicLinkDAO) FindByToken(_ string) (*models.MagicLink, error) {
	if m.FindErr != nil {
		return nil, m.FindErr
	}
	return m.FindResult, nil
}

func (m *MockMagicLinkDAO) MarkAsUsed(token string) error {
	m.MarkedAsUsed = append(m.MarkedAsUsed, token)
	return m.MarkAsUsedErr
}

type MockAuthTokenService struct {
	JWT string
	Err error
}

func (m *MockAuthTokenService) GenerateJWT(_ string) (string, error) {
	return m.JWT, m.Err
}

var ErrGeneric = errors.New("something went wrong")
