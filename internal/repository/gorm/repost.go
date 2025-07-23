package gorm

import (
	"context"
	"mini-blog/entity"
)

func (r *GormRepository) CreateRepost(ctx context.Context, repost entity.Repost) (*entity.Repost, error) {
	logger := r.Logger.With().Str("method", "CreateRepost").Logger()
	err := r.DB.WithContext(ctx).Create(&repost).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to create repost")
		return nil, err
	}

	return &repost, nil
}

func (r *GormRepository) GetReposts(ctx context.Context) ([]*entity.Repost, error) {
	logger := r.Logger.With().Str("method", "GetReposts").Logger()
	var reposts []*entity.Repost
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL").Find(&reposts).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get reposts")
		return nil, err
	}

	return reposts, nil
}

func (r *GormRepository) GetRepostByID(ctx context.Context, repostID int) (*entity.Repost, error) {
	logger := r.Logger.With().Str("method", "GetRepostByID").Int("repost_id", repostID).Logger()
	var repost entity.Repost
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL").First(&repost, repostID).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get repost by id")
		return nil, err
	}

	return &repost, nil
}

func (r *GormRepository) UpdateRepost(ctx context.Context, repost entity.Repost) (*entity.Repost, error) {
	logger := r.Logger.With().Str("method", "UpdateRepost").Int("repost_id", repost.ID).Logger()
	err := r.DB.WithContext(ctx).Model(&entity.Repost{}).Where("id = ?", repost.ID).Updates(&repost).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to update repost")
		return nil, err
	}

	return &repost, nil
}

func (r *GormRepository) DeleteRepost(ctx context.Context, repost entity.Repost) error {
	logger := r.Logger.With().Str("method", "DeleteRepost").Int("repost_id", repost.ID).Logger()
	err := r.DB.WithContext(ctx).Where("id = ?", repost.ID).Updates(&repost).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete repost")
		return err
	}

	return nil
}

func (r *GormRepository) GetRepostsByUserID(ctx context.Context, userID int) ([]*entity.Repost, error) {
	logger := r.Logger.With().Str("method", "GetRepostsByUserID").Int("user_id", userID).Logger()
	var reposts []*entity.Repost
	err := r.DB.WithContext(ctx).Where("user_id = ? AND deleted_at IS NULL", userID).Find(&reposts).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get reposts by user id")
		return nil, err
	}

	return reposts, nil
}
