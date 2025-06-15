package service

import (
	"mini-blog/internal/repository/gorm"

	"github.com/rs/zerolog"
)

type Service struct {
	Repository gorm.GormRepository
	Logger     zerolog.Logger
}

type ServiceI interface {
	// All service methods here
}
