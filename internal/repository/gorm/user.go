package gorm

import (
	"context"
	"mini-blog/entity"
)

func (r *GormRepository) CreateUser(ctx context.Context, u entity.User) (*entity.User, error) {
	logger := r.Logger.With().Str("method", "CreateUser").Logger()
	err := r.DB.WithContext(ctx).Create(&u).Error
	if err != nil {
		logger.Error().Err(err).Msg("error creating user")
		return nil, err
	}

	return &u, nil
}

func (r *GormRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	logger := r.Logger.With().Ctx(ctx).Str("method", "GetUserByEmail").Str("email", email).Logger()

	var user *entity.User
	err := r.DB.WithContext(ctx).Where("email = ? AND deleted_at IS NULL", email).First(&user).Error
	if err != nil {
		logger.Error().Err(err).Msg("error finding user by email")
		return nil, err
	}

	return user, nil
}

func (r *GormRepository) GetUserByID(ctx context.Context, userID int) (*entity.User, error) {
	logger := r.Logger.With().Ctx(ctx).Str("method", "GetUserByID").Int("user_id", userID).Logger()

	var user *entity.User
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL").First(&user, userID).Error
	if err != nil {
		logger.Error().Err(err).Msg("error finding user by id")
		return nil, err
	}

	return user, nil
}
