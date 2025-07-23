package main

import (
	"context"
	"mini-blog/config"
	"mini-blog/internal/controller"
	authv1 "mini-blog/internal/controller/auth/v1"
	"os"
	"os/signal"
	"syscall"
	"time"

	restv1 "mini-blog/internal/controller/rest/v1"
	"mini-blog/internal/repository/gorm"
	"mini-blog/internal/service"
	"net/http"

	"mini-blog/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}
	log.Info().Msg("start app")

	if _, err := config.InitConfig(); err != nil {
		log.Error().Err(err).Msg("failed int config")
		return
	}

	db, err := gorm.DBConnection(config.GetEnv().Db.Dsn, log)
	if err != nil {
		log.Error().Err(err).Msg("failed connection to database")
		return
	}

	repo := gorm.NewGormRepository(db, log)
	svc := service.NewService(repo, log)
	
	// - создается экземпляр gin Engine и к нему подключается миддвар CORS
	//  чтобы разрешать запросы с разных источников (например, с фронтенда)
	api := gin.New()
	api.Use(controller.CORS())

	authApi := authv1.AuthHandlerV1{
		Service: svc,
		Logger:  log,
	}
	authApi.InitRoutes(api)

	restApi := restv1.RestHandlerv1{
		Service: svc,
		Logger:  log,
	}
	restApi.InitRoutes(api)

	httpServer := http.Server{
		Addr:    config.Get().App.Port,
		Handler: api,
	}
 
	// Запускаем сервер в отдельной горутине
	go func() {
		log.Info().Msgf("Starting server on %s", config.Get().App.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("failed running server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("failed shutdown server")
	}

	log.Info().Msg("Server stopped")
}
