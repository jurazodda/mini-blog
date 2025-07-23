package service

import (
	"mini-blog/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndParseToken_Success(t *testing.T) {
	t.Setenv("JWT_SIGNINGKEY", "test-signing-key")
	t.Setenv("SIGNING_KEY", "test-signing-key")
	t.Setenv("TOKEN_TTL", "1")
	config.Set(&config.Config{
		Jwt: config.Jwt{
			TokenTTL:   1,
		},
	})
	token, err := GenerateToken(nil, 42)
	assert.NoError(t, err)
	claims, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, float64(42), claims["user_id"])
}

func TestParseToken_Invalid(t *testing.T) {
	_, err := ParseToken("invalid.token.here")
	assert.Error(t, err)
}
