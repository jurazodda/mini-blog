package gorm

import (
	"context"
	"mini-blog/entity"
)

func (r *GormRepository) CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	logger := r.Logger.With().Str("method", "CreatePost").Logger()
	err := r.DB.WithContext(ctx).Create(&post).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to create post")
		return nil, err
	}

	return &post, nil
}

func (r *GormRepository) GetPosts(ctx context.Context, userID int) ([]*entity.Post, error) {
	logger := r.Logger.With().Str("method", "GetPosts").Logger()
	var posts []*entity.Post
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL AND (user_id = ? OR is_public = true)", userID).Find(&posts).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get posts")
		return nil, err
	}

	return posts, nil
}

func (r *GormRepository) SearchPosts(ctx context.Context, userID int, query string, offset, limit int) ([]*entity.Post, int64, error) {
	logger := r.Logger.With().Str("method", "SearchPosts").Str("query", query).Logger()

	var posts []*entity.Post
	var total int64

	// Build the search query
	searchQuery := "%" + query + "%"
	baseQuery := r.DB.WithContext(ctx).Where("deleted_at IS NULL AND (user_id = ? OR is_public = true) AND (title ILIKE ? OR content ILIKE ?)", userID, searchQuery, searchQuery)

	// Count total results
	err := baseQuery.Model(&entity.Post{}).Count(&total).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to count search results")
		return nil, 0, err
	}

	// Get paginated results
	err = baseQuery.Offset(offset).Limit(limit).Order("created_at DESC").Find(&posts).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to search posts")
		return nil, 0, err
	}

	return posts, total, nil
}

func (r *GormRepository) GetPostByID(ctx context.Context, postID int) (*entity.Post, error) {
	logger := r.Logger.With().Str("method", "GetPostByID").Int("post_id", postID).Logger()
	var post entity.Post
	err := r.DB.WithContext(ctx).Where("deleted_at IS NULL").First(&post, postID).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to get post by id")
		return nil, err
	}

	return &post, nil
}

func (r *GormRepository) UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error) {
	logger := r.Logger.With().Str("method", "UpdatePost").Int("post_id", post.ID).Logger()
	err := r.DB.WithContext(ctx).Model(&entity.Post{}).Where("id = ?", post.ID).Updates(&post).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to update post")
		return nil, err
	}

	return &post, nil
}

func (r *GormRepository) DeletePost(ctx context.Context, post entity.Post) error {
	logger := r.Logger.With().Str("method", "DeletePost").Int("post_id", post.ID).Logger()
	err := r.DB.WithContext(ctx).Where("id = ?", post.ID).Updates(&post).Error
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete post")
		return err
	}

	return nil
}
