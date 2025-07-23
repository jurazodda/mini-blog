package service

import (
	"context"
	"mini-blog/entity"
	"time"
)

func (s *Service) CreateComment(ctx context.Context, in CreateCommentIn, userID int) (*entity.Comment, error) {
	logger := s.Logger.With().Str("method", "CreateComment").Int("user_id", userID).Logger()
	if err := in.Validate(); err != nil {
		logger.Error().Err(err).Msg("failed to validate comment")
		return nil, err
	}

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	createdCommmnet, err := s.Repository.CreateComment(ctx, entity.Comment{
		UserID:  userID,
		PostID:  in.PostID,
		Comment: in.Comment,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create comment")
		return nil, err
	}

	return createdCommmnet, nil
}

func (s *Service) GetCommentByID(ctx context.Context, userID, commentID int) (*entity.Comment, error) {
	logger := s.Logger.With().Str("method", "GetCommentByID").Int("user_id", userID).Int("comment_id", commentID).Logger()
	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	if commentID <= 0 {
		logger.Error().Msg("comment id is required")
		return nil, ErrCommentIdRequired
	}

	comment, err := s.Repository.GetCommentByID(ctx, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comment by id")
		return nil, err
	}

	return comment, nil
}

func (s *Service) UpdateComment(ctx context.Context, comment string, userID, commentID int) (*entity.Comment, error) {
	logger := s.Logger.With().Str("method", "UpdateComment").Int("user_id", userID).Int("comment_id", commentID).Logger()

	if comment == "" {
		logger.Error().Msg("comment is required")
		return nil, ErrCommentRequired
	}

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	if commentID <= 0 {
		logger.Error().Msg("comment id is required")
		return nil, ErrCommentIdRequired
	}

	// Проверяем, что комментарий принадлежит пользователю
	existingComment, err := s.Repository.GetCommentByID(ctx, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get existing comment")
		return nil, err
	}

	if existingComment.UserID != userID {
		logger.Error().Msg("user can only update own comments")
		return nil, ErrUnauthorized
	}

	updatedComment, err := s.Repository.UpdateComment(ctx, entity.Comment{
		ID:      commentID,
		UserID:  userID,
		Comment: comment,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update comment")
		return nil, err
	}

	return updatedComment, nil
}

func (s *Service) DeleteComment(ctx context.Context, userID, commentID int) error {
	logger := s.Logger.With().Str("method", "DeleteComment").Int("user_id", userID).Int("comment_id", commentID).Logger()

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return ErrUserIdRequired
	}

	if commentID <= 0 {
		logger.Error().Msg("comment id is required")
		return ErrCommentIdRequired
	}

	// Проверяем, что комментарий принадлежит пользователю
	existingComment, err := s.Repository.GetCommentByID(ctx, commentID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get existing comment")
		return err
	}

	if existingComment.UserID != userID {
		logger.Error().Msg("user can only delete own comments")
		return ErrUnauthorized
	}

	tNow := time.Now()
	toUpdate := entity.Comment{
		ID:        commentID,
		DeletedAt: &tNow,
	}

	err = s.Repository.DeleteComment(ctx, toUpdate)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete comment")
		return err
	}

	logger.Info().Msg("comment deleted successfully")
	return nil
}

func (s *Service) GetCommentsByPostID(ctx context.Context, postID int) ([]*entity.Comment, error) {
	logger := s.Logger.With().Str("method", "GetCommentsByPostID").Int("post_id", postID).Logger()

	if postID <= 0 {
		logger.Error().Msg("post id is required")
		return nil, ErrPostIdRequired
	}

	comments, err := s.Repository.GetCommentsByPostID(ctx, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comments by post id")
		return nil, err
	}

	return comments, nil
}
