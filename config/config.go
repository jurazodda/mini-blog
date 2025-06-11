package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	config *Config
)

type Config struct {
	App App
	DB  DB
}

type App struct {
	Port string
}

type DB struct {
	Dsn string
}

func InitConfig() (*Config, error) {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("config")

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error().Err(err).Msg("error reading config file")
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config = &Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		logger.Error().Err(err).Msg("unable to decode into struct")
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		logger.Info().Msgf("Config file changed: %s", in.Name)
		err = viper.Unmarshal(&config)
		if err != nil {
			logger.Error().Err(err).Msg("unable to decode into struct")
		}
	})

	return config, nil
}
