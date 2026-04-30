package bo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magic-link-auth/src/layers/main/bo"
	"magic-link-auth/tests/testutil"
)

func TestCreateMagicLinkBO_Success(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{}
	email := &testutil.MockEmailService{}
	token := &testutil.MockTokenService{Token: "abc123"}

	b := bo.NewCreateMagicLinkBO(dao, email, token, "http://localhost:8080")
	err := b.Execute("user@example.com")

	require.NoError(t, err)
	assert.Equal(t, []string{"user@example.com"}, email.CalledWith)
}

func TestCreateMagicLinkBO_TokenGenerateError(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{}
	email := &testutil.MockEmailService{}
	token := &testutil.MockTokenService{Err: testutil.ErrGeneric}

	b := bo.NewCreateMagicLinkBO(dao, email, token, "http://localhost:8080")
	err := b.Execute("user@example.com")

	require.ErrorIs(t, err, testutil.ErrGeneric)
	assert.Empty(t, email.CalledWith)
}

func TestCreateMagicLinkBO_SaveError(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{SaveErr: testutil.ErrGeneric}
	email := &testutil.MockEmailService{}
	token := &testutil.MockTokenService{Token: "abc123"}

	b := bo.NewCreateMagicLinkBO(dao, email, token, "http://localhost:8080")
	err := b.Execute("user@example.com")

	require.ErrorIs(t, err, testutil.ErrGeneric)
	assert.Empty(t, email.CalledWith)
}

func TestCreateMagicLinkBO_SendEmailError(t *testing.T) {
	dao := &testutil.MockMagicLinkDAO{}
	email := &testutil.MockEmailService{Err: testutil.ErrGeneric}
	token := &testutil.MockTokenService{Token: "abc123"}

	b := bo.NewCreateMagicLinkBO(dao, email, token, "http://localhost:8080")
	err := b.Execute("user@example.com")

	require.ErrorIs(t, err, testutil.ErrGeneric)
}
