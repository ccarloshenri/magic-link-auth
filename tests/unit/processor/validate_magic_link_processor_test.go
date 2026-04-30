package processor_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magic-link-auth/src/layers/main/bo"
	"magic-link-auth/src/layers/main/enums"
	"magic-link-auth/src/layers/main/models"
	"magic-link-auth/src/layers/main/processor"
	"magic-link-auth/tests/testutil"
)

func newValidateProcessor(link *models.MagicLink, findErr error, jwt string, jwtErr error) *processor.ValidateMagicLinkProcessor {
	dao := &testutil.MockMagicLinkDAO{FindResult: link, FindErr: findErr}
	auth := &testutil.MockAuthTokenService{JWT: jwt, Err: jwtErr}
	b := bo.NewValidateMagicLinkBO(dao, auth)
	return processor.NewValidateMagicLinkProcessor(b)
}

func TestValidateMagicLinkProcessor_EmptyToken(t *testing.T) {
	p := newValidateProcessor(nil, nil, "", nil)
	_, err := p.Process("")
	assert.EqualError(t, err, "token is required")
}

func TestValidateMagicLinkProcessor_Success(t *testing.T) {
	link := &models.MagicLink{
		Token:     "tok",
		Email:     "user@example.com",
		ExpiresAt: time.Now().Add(10 * time.Minute).Unix(),
		Status:    enums.Pending,
	}
	p := newValidateProcessor(link, nil, "jwt-token", nil)
	out, err := p.Process("tok")

	require.NoError(t, err)
	assert.Equal(t, "jwt-token", out.AccessToken)
	assert.Equal(t, "Bearer", out.Type)
}

func TestValidateMagicLinkProcessor_TokenNotFound(t *testing.T) {
	p := newValidateProcessor(nil, testutil.ErrGeneric, "", nil)
	_, err := p.Process("tok")
	assert.ErrorIs(t, err, bo.ErrTokenNotFound)
}
