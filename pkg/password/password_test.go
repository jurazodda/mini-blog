package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePasswordHashAndCompare(t *testing.T) {
	// Test-case: проверяем, что GeneratePasswordHash возвращает не nil
	password := "supersecret123"
	hash, err := GeneratePasswordHash(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	// Test-case: проверяем, что ComparePassword возвращает true для правильного пароля
	ok := ComparePassword(hash, password)
	assert.True(t, ok)

	// Test-case: проверяем, что ComparePassword возвращает false для неправильного пароля
	ok = ComparePassword(hash, "wrongpassword")
	assert.False(t, ok)
}

func TestGeneratePasswordHash_Error(t *testing.T) {
	// Test-case: проверяем, что GeneratePasswordHash возвращает не nil
	hash, err := GeneratePasswordHash("")
	
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
}
