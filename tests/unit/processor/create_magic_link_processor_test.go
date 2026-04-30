package processor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"magic-link-auth/src/layers/main/bo"
	"magic-link-auth/src/layers/main/processor"
	"magic-link-auth/tests/testutil"
)

func newCreateProcessor(token string, tokenErr, saveErr, emailErr error) *processor.CreateMagicLinkProcessor {
	dao := &testutil.MockMagicLinkDAO{SaveErr: saveErr}
	email := &testutil.MockEmailService{Err: emailErr}
	ts := &testutil.MockTokenService{Token: token, Err: tokenErr}
	b := bo.NewCreateMagicLinkBO(dao, email, ts, "http://localhost:8080")
	return processor.NewCreateMagicLinkProcessor(b)
}

func TestCreateMagicLinkProcessor_EmptyEmail(t *testing.T) {
	p := newCreateProcessor("tok", nil, nil, nil)
	err := p.Process(processor.CreateMagicLinkInput{Email: ""})
	assert.EqualError(t, err, "email is required")
}

func TestCreateMagicLinkProcessor_InvalidEmail(t *testing.T) {
	p := newCreateProcessor("tok", nil, nil, nil)
	err := p.Process(processor.CreateMagicLinkInput{Email: "not-an-email"})
	assert.ErrorContains(t, err, "invalid email")
}

func TestCreateMagicLinkProcessor_Success(t *testing.T) {
	p := newCreateProcessor("tok", nil, nil, nil)
	err := p.Process(processor.CreateMagicLinkInput{Email: "user@example.com"})
	require.NoError(t, err)
}

func TestCreateMagicLinkProcessor_BOError(t *testing.T) {
	p := newCreateProcessor("", testutil.ErrGeneric, nil, nil)
	err := p.Process(processor.CreateMagicLinkInput{Email: "user@example.com"})
	require.ErrorIs(t, err, testutil.ErrGeneric)
}
