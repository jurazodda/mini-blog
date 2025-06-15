package main

import (
	"context"
	"mini-blog/config"
	authv1 "mini-blog/internal/controller/auth/v1"
	"time"

	restv1 "mini-blog/internal/controller/rest/v1"
	"mini-blog/internal/repository/gorm"
	"mini-blog/internal/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

func main() {
	ctx := context.Background()
	logger := zerolog.New(os.Stdout).With().Caller().Logger()
	logger.Info().Msg("start app")
	_, err := config.InitConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed int config")
	}

	db, err := gorm.DBConnection(config.Get().DB.Dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed connection to database")
	}

	repo := gorm.GormRepository{
		DB:     db,
		Logger: logger.With().Str("component", "Repository").Logger(),
	}

	svc := service.Service{
		Repository: repo,
		Logger:     logger.With().Str("component", "Service").Logger(),
	}

	api := gin.New()

	authApi := authv1.AuthHandlerV1{
		Service: svc,
		Logger:  logger.With().Str("component", "AuthHandlerV1").Logger(),
	}
	authApi.InitRoutes(api)

	restApi := restv1.RestHandlerv1{
		Service: svc,
		Logger:  logger.With().Str("component", "RestHandlerv1").Logger(),
	}
	restApi.InitRoutes(api)

	httpServer := http.Server{
		Addr:    config.Get().App.Port,
		Handler: api,
	}

	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed running server")
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err = httpServer.Shutdown(shutdownCtx)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed shutdown server")
	}
}
