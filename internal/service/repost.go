package service

import (
	"context"
	"mini-blog/entity"
	"time"
)

func (s *Service) CreateRepost(ctx context.Context, in CreateRepostIn, userID int) (*entity.Repost, error) {
	logger := s.Logger.With().Str("method", "CreateRepost").Int("user_id", userID).Int("post_id", in.PostID).Logger()

	if err := in.Validate(); err != nil {
		logger.Error().Err(err).Msg("failed to validate create repost")
		return nil, err
	}

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	// Проверяем, существует ли оригинальный пост
	originalPost, err := s.Repository.GetPostByID(ctx, in.PostID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get original post")
		return nil, err
	}

	// Создаем репост
	repost := entity.Repost{
		UserID:   userID,
		PostID:   in.PostID,
		Title:    in.Title,
		Content:  in.Content,
		ImageURL: &originalPost.ImageURL,
		Comments: 0,
		Likes:    0,
	}

	// Если пользователь не указал свой текст, используем оригинальный
	if repost.Title == "" {
		repost.Title = originalPost.Title
	}
	if repost.Content == "" {
		repost.Content = originalPost.Content
	}

	createdRepost, err := s.Repository.CreateRepost(ctx, repost)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create repost")
		return nil, err
	}

	logger.Info().Msg("repost created successfully")
	return createdRepost, nil
}

func (s *Service) GetReposts(ctx context.Context) ([]*entity.Repost, error) {
	logger := s.Logger.With().Str("method", "GetReposts").Logger()

	reposts, err := s.Repository.GetReposts(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get reposts")
		return nil, err
	}

	return reposts, nil
}

func (s *Service) GetRepostByID(ctx context.Context, userID, repostID int) (*entity.Repost, error) {
	logger := s.Logger.With().Str("method", "GetRepostByID").Int("user_id", userID).Int("repost_id", repostID).Logger()

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	if repostID <= 0 {
		logger.Error().Msg("repost id is required")
		return nil, ErrPostIdRequired
	}

	repost, err := s.Repository.GetRepostByID(ctx, repostID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get repost by id")
		return nil, err
	}

	return repost, nil
}

func (s *Service) UpdateRepost(ctx context.Context, in UpdateRepostIn, userID, repostID int) (*entity.Repost, error) {
	logger := s.Logger.With().Str("method", "UpdateRepost").Int("user_id", userID).Int("repost_id", repostID).Logger()

	if err := in.Validate(); err != nil {
		logger.Error().Err(err).Msg("failed to validate update repost")
		return nil, err
	}

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	if repostID <= 0 {
		logger.Error().Msg("repost id is required")
		return nil, ErrPostIdRequired
	}

	// Проверяем, что репост принадлежит пользователю
	existingRepost, err := s.Repository.GetRepostByID(ctx, repostID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get existing repost")
		return nil, err
	}

	if existingRepost.UserID != userID {
		logger.Error().Msg("user can only update own reposts")
		return nil, ErrUnauthorized
	}

	updatedRepost, err := s.Repository.UpdateRepost(ctx, entity.Repost{
		ID:      repostID,
		UserID:  userID,
		Title:   in.Title,
		Content: in.Content,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update repost")
		return nil, err
	}

	logger.Info().Msg("repost updated successfully")
	return updatedRepost, nil
}

func (s *Service) DeleteRepost(ctx context.Context, userID, repostID int) error {
	logger := s.Logger.With().Str("method", "DeleteRepost").Int("user_id", userID).Int("repost_id", repostID).Logger()

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return ErrUserIdRequired
	}

	if repostID <= 0 {
		logger.Error().Msg("repost id is required")
		return ErrPostIdRequired
	}

	// Проверяем, что репост принадлежит пользователю
	existingRepost, err := s.Repository.GetRepostByID(ctx, repostID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get existing repost")
		return err
	}

	if existingRepost.UserID != userID {
		logger.Error().Msg("user can only delete own reposts")
		return ErrUnauthorized
	}

	tNow := time.Now()
	toUpdate := entity.Repost{
		ID:        repostID,
		DeletedAt: &tNow,
	}

	err = s.Repository.DeleteRepost(ctx, toUpdate)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete repost")
		return err
	}

	logger.Info().Msg("repost deleted successfully")
	return nil
}

func (s *Service) GetRepostsByUserID(ctx context.Context, targetUserID int) ([]*entity.Repost, error) {
	logger := s.Logger.With().Str("method", "GetRepostsByUserID").Int("target_user_id", targetUserID).Logger()

	if targetUserID <= 0 {
		logger.Error().Msg("target user id is required")
		return nil, ErrUserIdRequired
	}

	reposts, err := s.Repository.GetRepostsByUserID(ctx, targetUserID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get reposts by user id")
		return nil, err
	}

	return reposts, nil
}
