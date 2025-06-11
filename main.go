package main

import (
	"mini-blog/config"
	"os"

	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	_, err := config.InitConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed init config")
	}

	
}
