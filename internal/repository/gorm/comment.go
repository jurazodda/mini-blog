package gorm

import (
	"context"
	"mini-blog/entity"
)

func (r *GormRepository) CreateComment(ctx context.Context, comment entity.Comment) (*entity.Comment, error) {
	logger := r.Logger.With().Str("method", "CreateComment").Logger()
	err := r.DB.WithContext(ctx).Create(&comment).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to create comment")
		return nil, err
	}

	return &comment, nil
}

func (r *GormRepository) GetCommentsByPostID(ctx context.Context, postID int) ([]*entity.Comment, error) {
	logger := r.Logger.With().Str("method", "GetCommentsByPostID").Int("post_id", postID).Logger()
	var comments []*entity.Comment
	err := r.DB.WithContext(ctx).Where("post_id = ? AND deleted_at IS NULL", postID).Find(&comments).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comments by post id")
		return nil, err
	}

	return comments, nil
}

func (r *GormRepository) GetCommentsByCommentLikesDesc(ctx context.Context, postID int) ([]*entity.Comment, error) {
	logger := r.Logger.With().Str("method", "GetCommentsByCommentLikesDesc").Int("post_id", postID).Logger()
	comments := make([]*entity.Comment, 0)
	err := r.DB.WithContext(ctx).
		Where("post_id = ? AND deleted_at IS NULL", postID).
		Find(&comments).
		Order("comment_like desc").Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comments by comment likes desc")
		return nil, err
	}
	return comments, nil
}

func (r *GormRepository) CountComments(ctx context.Context, postID int) (int, error) {
	logger := r.Logger.With().Str("method", "CountComments").Int("post_id", postID).Logger()
	var comment entity.Comment
	var count int64
	err := r.DB.WithContext(ctx).Model(&comment).Where("post_id = ? AND deleted_at IS NULL", postID).Count(&count).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to count comments")
		return 0, err
	}

	return int(count), nil
}

func (r *GormRepository) GetCommentByID(ctx context.Context, commentID int) (*entity.Comment, error) {
	logger := r.Logger.With().Str("method", "GetCommentByID").Int("comment_id", commentID).Logger()
	var comment entity.Comment
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL").First(&comment, commentID).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comment by id")
		return nil, err
	}

	return &comment, nil
}

func (r *GormRepository) UpdateComment(ctx context.Context, comment entity.Comment) (*entity.Comment, error) {
	logger := r.Logger.With().Str("method", "UpdateComment").Int("comment_id", comment.ID).Logger()
	err := r.DB.WithContext(ctx).Model(&entity.Comment{}).Where("id = ?", comment.ID).Updates(&comment).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to update comment")
		return nil, err
	}

	return &comment, nil
}

func (r *GormRepository) DeleteComment(ctx context.Context, comment entity.Comment) error {
	logger := r.Logger.With().Str("method", "DeleteComment").Int("comment_id", comment.ID).Logger()
	err := r.DB.WithContext(ctx).Where("id = ?", comment.ID).Updates(&comment).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete comment")
		return err
	}

	return nil
}
