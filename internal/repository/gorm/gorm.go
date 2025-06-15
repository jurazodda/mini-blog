package gorm

import (
	"log"
	"mini-blog/internal/entity"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func DBConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("error connection to database: ", err.Error())
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Post{}, &entity.Post{}, &entity.Comment{}, &entity.Like{}, &entity.Repost{})
	if err != nil {
		log.Fatal("error migrating database")
	}

	return db, nil
}
