package gorm

import (
	"context"
	"mini-blog/internal/entity"
)

func (r *GormRepository) CreateUser(ctx context.Context, u entity.User) error {
	err := r.DB.WithContext(ctx).Create(&u).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GormRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	logger := r.Logger.With().Ctx(ctx).Str("method", "GetUserByEmail").Str("email", email).Logger()

	var user *entity.User
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		logger.Error().Err(err).Msg("error finding user by email")
		return nil, translateError(err)
	}

	return user, nil
}