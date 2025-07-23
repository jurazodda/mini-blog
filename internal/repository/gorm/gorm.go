package gorm

import (
	"mini-blog/entity"
	"mini-blog/internal/repository"

	"mini-blog/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormRepository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

// Ensure GormRepository implements Repository interface
var _ repository.RepositoryI = (*GormRepository)(nil)

func DBConnection(dsn string, log *logger.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("error connection to database")
		return nil, err
	}

	log.Info().Msg("Starting GORM database migrations")
	err = db.AutoMigrate(&entity.User{}, &entity.Post{}, &entity.Comment{}, &entity.Like{}, &entity.CommentLike{}, &entity.Repost{})
	if err != nil {
		log.Error().Err(err).Msg("Failed to migrate database")
		return nil, err
	}

	return db, nil
}

func NewGormRepository(db *gorm.DB, logger *logger.Logger) *GormRepository {
	return &GormRepository{
		DB:     db,
		Logger: logger,
	}
}
