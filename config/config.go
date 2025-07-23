package config

import (
	"fmt"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	config *Config
)

type Config struct {
	App App
	Jwt Jwt
}

type App struct {
	Port string
}

type Jwt struct {
	TokenTTL   time.Duration
}

func Get() Config {
	return *config
}

func InitConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	config = &Config{}
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		log.Println("Config file changed", in.Name)
		if err := viper.Unmarshal(&config); err != nil {
			log.Printf("unable to decode into struct: %v", err)
		}
	})

	return config, nil
}

func Set(cfg *Config) {
	config = cfg
}
