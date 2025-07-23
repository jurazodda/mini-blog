package gorm

import (
	"context"
	"mini-blog/entity"
)

func (r *GormRepository) GetLikesByPostID(ctx context.Context, postID int) ([]*entity.Like, error) {
	logger := r.Logger.With().Str("method", "GetLikesByPostID").Int("post_id", postID).Logger()
	var likes []*entity.Like
	err := r.DB.WithContext(ctx).Where("post_id = ? AND deleted_at IS NULL", postID).Find(&likes).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get likes by post id")
		return nil, err
	}

	return likes, nil
}

func (r *GormRepository) CountLikes(ctx context.Context, postID int) (int, error) {
	logger := r.Logger.With().Str("method", "CountLikes").Int("post_id", postID).Logger()
	var like entity.Like
	var count int64
	err := r.DB.WithContext(ctx).Model(&like).Where("post_id = ? AND deleted_at IS NULL", postID).Count(&count).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to count likes")
		return 0, err
	}

	return int(count), nil
}

func (r *GormRepository) GetLikeByUserAndPost(ctx context.Context, userID, postID int) (*entity.Like, error) {
	logger := r.Logger.With().Str("method", "GetLikeByUserAndPost").Int("user_id", userID).Int("post_id", postID).Logger()
	var like entity.Like
	err := r.DB.WithContext(ctx).Where("user_id = ? AND post_id = ? AND deleted_at IS NULL", userID, postID).First(&like).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get like by user and post")
		return nil, err
	}

	return &like, nil
}

func (r *GormRepository) CreateLike(ctx context.Context, l entity.Like) (*entity.Like, error) {
	logger := r.Logger.With().Str("method", "CreateLike").Logger()
	err := r.DB.WithContext(ctx).Create(&l).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to create like")
		return nil, err
	}

	return &l, nil
}

func (r *GormRepository) DeleteLike(ctx context.Context, l entity.Like) error {
	logger := r.Logger.With().Str("method", "DeleteLike").Logger()
	err := r.DB.WithContext(ctx).Where("user_id = ? AND post_id = ? AND deleted_at IS NULL", l.UserID, l.PostID).Delete(&l).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete like")
		return err
	}

	return nil
}

func (r *GormRepository) GetLikeByUserAndComment(ctx context.Context, userID, commentID int) (*entity.CommentLike, error) {
	logger := r.Logger.With().Str("method", "GetLikeByUserAndComment").Int("user_id", userID).Int("comment_id", commentID).Logger()
	var like entity.CommentLike
	err := r.DB.WithContext(ctx).Where("user_id = ? AND comment_id = ?", userID, commentID).First(&like).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get like by user and comment")
		return nil, err
	}

	return &like, nil
}

func (r *GormRepository) CreateLikeComment(ctx context.Context, like entity.CommentLike) (*entity.CommentLike, error) {
	logger := r.Logger.With().Str("method", "CreateLikeComment").Logger()
	err := r.DB.WithContext(ctx).Create(&like).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to create comment like")
		return nil, err
	}

	return &like, nil
}

func (r *GormRepository) DeleteLikeComment(ctx context.Context, likeID int) error {
	logger := r.Logger.With().Str("method", "DeleteLikeComment").Int("like_id", likeID).Logger()
	err := r.DB.WithContext(ctx).Delete(&entity.CommentLike{}, likeID).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete comment like")
		return err
	}

	return nil
}
