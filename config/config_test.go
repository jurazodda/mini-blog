package config

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetAndGet(t *testing.T) {
	// Test-case: проверяем, что можно установить и получить конфиг
	cfg := &Config{
		App: App{Port: "8080"},
	}

	Set(cfg)
	got := Get()

	assert.Equal(t, cfg.App.Port, got.App.Port)
	assert.Equal(t, cfg.Jwt.TokenTTL, got.Jwt.TokenTTL)
}

func TestInitConfig_Success(t *testing.T) {
	// Test-case: проверяем, что можно инициализировать конфиг из файла
	dir := t.TempDir()
	file := dir + "/config.yml"
	content := []byte("db:\n  dsn: test-dsn\napp:\n  port: '8080'\njwt:\n  tokenTTL: 12h\n  signingKey: test-key\n")
	os.WriteFile(file, content, 0644)

	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(dir)

	cfg, err := InitConfig()

	assert.NoError(t, err)
	assert.Equal(t, "8080", cfg.App.Port)
	assert.Equal(t, 12*time.Hour, cfg.Jwt.TokenTTL)
}

func TestInitConfig_FileNotFound(t *testing.T) {
	// Test-case: проверяем, что InitConfig возвращает ошибку, если файл не найден
	dir := t.TempDir()
	viper.Reset()
	viper.SetConfigName("notfound")
	viper.SetConfigType("yml")
	viper.AddConfigPath(dir)

	cfg, err := InitConfig()
	
	assert.Error(t, err)
	assert.Nil(t, cfg)
}

func TestGetEnv(t *testing.T) {
	// Test-case: проверяем, что можно получить переменные окружения
	t.Setenv("DSN_URL", "env-dsn")
	t.Setenv("SIGNING_KEY", "env-key")

	env := GetEnv()
	
	assert.Equal(t, "env-dsn", env.Db.Dsn)
	assert.Equal(t, "env-key", env.Jwt.SigningKey)
}
