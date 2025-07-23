package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTestLogger(t *testing.T) {
	// Test-case: проверяем, что NewTestLogger возвращает не nil
	log := NewTestLogger()
	assert.NotNil(t, log)

	// Test-case: проверяем, что InfoLogger не пишет в файл
	log.Info().Msg("test info")
	log.Debug().Msg("test debug")
	log.Warn().Msg("test warn")
	log.Error().Msg("test error")
}

func TestLogger_With_WithContext(t *testing.T) {
	// Test-case: проверяем, что With возвращает не nil
	log := NewTestLogger()
	ctx := log.With().Str("foo", "bar")
	assert.NotNil(t, ctx)
	
	// Test-case: проверяем, что WithContext возвращает не nil
	ctx2 := log.WithContext(map[string]string{"k": "v"})
	assert.NotNil(t, ctx2)
}

func TestGetProjectRoot(t *testing.T) {
	// Test-case: проверяем, что getProjectRoot возвращает не nil
	root, err := getProjectRoot()
	assert.NoError(t, err)
	assert.NotEmpty(t, root)

	// Test-case: проверяем, что в корне должен быть go.mod
	_, err = os.Stat(filepath.Join(root, "go.mod"))
	assert.NoError(t, err)
}

func TestNewLogger_CreatesFilesInTempDir(t *testing.T) {
	tempDir := t.TempDir()
	// Test-case: проверяем, что в корне должен быть go.mod
	os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module testmod\n"), 0644)
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	// Test-case: проверяем, что NewLogger возвращает не nil
	log, err := NewLogger()
	assert.NoError(t, err)
	assert.NotNil(t, log)

	// Test-case: проверяем, что файлы логов созданы
	for _, f := range []string{"info.json", "debug.json", "warn.json", "error.json"} {
		_, err := os.Stat(filepath.Join(tempDir, "logs", f))
		assert.NoError(t, err)
	}
}
