package bo_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magic-link-auth/src/layers/main/bo"
	"magic-link-auth/src/layers/main/enums"
	"magic-link-auth/src/layers/main/models"
	"magic-link-auth/tests/testutil"
)

func pendingLink(email string, expiresAt int64) *models.MagicLink {
	return &models.MagicLink{
		Token:     "tok",
		Email:     email,
		ExpiresAt: expiresAt,
		Status:    enums.Pending,
	}
}

func TestValidateMagicLinkBO_Success(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{
		FindResult: pendingLink("user@example.com", time.Now().Add(10*time.Minute).Unix()),
	}
	auth := &testutil.MockAuthTokenService{JWT: "jwt-token"}

	b := bo.NewValidateMagicLinkBO(dao, auth)
	token, err := b.Execute("tok")

	require.NoError(t, err)
	assert.Equal(t, "jwt-token", token)
	assert.Equal(t, []string{"tok"}, dao.MarkedAsUsed)
}

func TestValidateMagicLinkBO_NotFound(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{FindErr: testutil.ErrGeneric}
	auth := &testutil.MockAuthTokenService{}

	b := bo.NewValidateMagicLinkBO(dao, auth)
	_, err := b.Execute("tok")

	assert.ErrorIs(t, err, bo.ErrTokenNotFound)
}

func TestValidateMagicLinkBO_AlreadyUsed(t *testing.T) {
	link := pendingLink("user@example.com", time.Now().Add(10*time.Minute).Unix())
	link.Status = enums.Used
	dao := &testutil.MockMagicLinkDAO{FindResult: link}
	auth := &testutil.MockAuthTokenService{}

	b := bo.NewValidateMagicLinkBO(dao, auth)
	_, err := b.Execute("tok")

	assert.ErrorIs(t, err, bo.ErrTokenUsed)
}

func TestValidateMagicLinkBO_StatusExpired(t *testing.T) {
	link := pendingLink("user@example.com", time.Now().Add(10*time.Minute).Unix())
	link.Status = enums.Expired
	dao := &testutil.MockMagicLinkDAO{FindResult: link}
	auth := &testutil.MockAuthTokenService{}

	b := bo.NewValidateMagicLinkBO(dao, auth)
	_, err := b.Execute("tok")

	assert.ErrorIs(t, err, bo.ErrTokenExpired)
}

func TestValidateMagicLinkBO_TimeExpired(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{
		FindResult: pendingLink("user@example.com", time.Now().Add(-1*time.Minute).Unix()),
	}
	auth := &testutil.MockAuthTokenService{}

	b := bo.NewValidateMagicLinkBO(dao, auth)
	_, err := b.Execute("tok")

	assert.ErrorIs(t, err, bo.ErrTokenExpired)
}

func TestValidateMagicLinkBO_MarkAsUsedError(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{
		FindResult:    pendingLink("user@example.com", time.Now().Add(10*time.Minute).Unix()),
		MarkAsUsedErr: testutil.ErrGeneric,
	}
	auth := &testutil.MockAuthTokenService{}

	b := bo.NewValidateMagicLinkBO(dao, auth)
	_, err := b.Execute("tok")

	require.ErrorIs(t, err, testutil.ErrGeneric)
}

func TestValidateMagicLinkBO_JWTError(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{
		FindResult: pendingLink("user@example.com", time.Now().Add(10*time.Minute).Unix()),
	}
	auth := &testutil.MockAuthTokenService{Err: testutil.ErrGeneric}

	b := bo.NewValidateMagicLinkBO(dao, auth)
	_, err := b.Execute("tok")

	require.ErrorIs(t, err, testutil.ErrGeneric)
}
