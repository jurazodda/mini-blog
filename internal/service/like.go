package service

import (
	"context"
	"errors"
	"mini-blog/entity"

	"gorm.io/gorm"
)

// ToggleLikeResult представляет результат операции toggle лайка
type ToggleLikeResult struct {
	IsLiked bool   `json:"is_liked"`
	Message string `json:"message"`
}

func (s *Service) CreateLikePost(ctx context.Context, postID, userID int) (*ToggleLikeResult, error) {
	logger := s.Logger.With().Str("method", "CreateLikePost").Int("user_id", userID).Int("post_id", postID).Logger()

	if postID <= 0 {
		logger.Error().Msg("post id is required")
		return nil, ErrPostIdRequired
	}

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	// Проверяем, существует ли пост
	_, err := s.Repository.GetPostByID(ctx, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get post by id")
		return nil, err
	}

	// Проверяем, есть ли уже лайк от этого пользователя
	existingLike, err := s.Repository.GetLikeByUserAndPost(ctx, userID, postID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error().Err(err).Msg("failed to check existing like")
		return nil, err
	}

	if existingLike != nil {
		// Лайк существует - удаляем его
		err := s.Repository.DeleteLike(ctx, entity.Like{
			UserID: userID,
			PostID: postID,
		})
		if err != nil {
			logger.Error().Err(err).Msg("failed to delete like")
			return nil, err
		}
		logger.Info().Msg("like removed")
		return &ToggleLikeResult{
			IsLiked: false,
			Message: "like removed",
		}, nil
	}

	// Лайка нет - создаем новый
	_, err = s.Repository.CreateLike(ctx, entity.Like{
		UserID: userID,
		PostID: postID,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create like")
		return nil, err
	}

	logger.Info().Msg("like created")
	return &ToggleLikeResult{
		IsLiked: true,
		Message: "like created",
	}, nil
}

func (s *Service) CreateLikeComment(ctx context.Context, commentID, userID int) (*ToggleLikeResult, error) {
	logger := s.Logger.With().Str("method", "CreateLikeComment").Int("user_id", userID).Int("comment_id", commentID).Logger()

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	if commentID <= 0 {
		logger.Error().Msg("comment id is required")
		return nil, ErrCommentIdRequired
	}

	// Проверяем, существует ли комментарий
	_, err := s.Repository.GetCommentByID(ctx, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comment")
		return nil, err
	}

	// Проверяем, есть ли уже лайк от этого пользователя
	existingLike, err := s.Repository.GetLikeByUserAndComment(ctx, userID, commentID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error().Err(err).Msg("failed to check existing like")
		return nil, err
	}

	if existingLike != nil {
		// Удаляем лайк (toggle off)
		err = s.Repository.DeleteLikeComment(ctx, existingLike.ID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to delete like")
			return nil, err
		}

		logger.Info().Msg("like removed from comment")
		return &ToggleLikeResult{
			IsLiked: false,
			Message: "like removed",
		}, nil
	} else {
		// Создаем новый лайк
		like := entity.CommentLike{
			UserID:    userID,
			CommentID: commentID,
		}

		_, err = s.Repository.CreateLikeComment(ctx, like)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create like")
			return nil, err
		}

		logger.Info().Msg("like added to comment")
		return &ToggleLikeResult{
			IsLiked: true,
			Message: "like created",
		}, nil
	}
}
