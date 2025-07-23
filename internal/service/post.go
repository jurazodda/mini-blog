package service

import (
	"context"
	"mini-blog/entity"
	"time"
)

func (s *Service) CreatePost(ctx context.Context, in CreatePostIn, userID int) (*entity.Post, error) {
	logger := s.Logger.With().Str("method", "CreatePost").Int("user_id", userID).Logger()
	if err := in.Validate(); err != nil {
		logger.Error().Err(err).Msg("failed to validate create post")
		return nil, err
	}

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	createdPost, err := s.Repository.CreatePost(ctx, entity.Post{
		UserID:   userID,
		Title:    in.Title,
		ImageURL: in.ImageURL,
		Content:  in.Content,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to create post")
		return nil, err
	}

	return createdPost, nil
}

func (s *Service) GetPosts(ctx context.Context, userID int) ([]*entity.PostList, error) {
	logger := s.Logger.With().Str("method", "GetPosts").Logger()
	posts, err := s.Repository.GetPosts(ctx, userID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get posts")
		return nil, err
	}

	var postList []*entity.PostList
	for _, post := range posts {
		commtentsCount, err := s.Repository.CountComments(ctx, post.ID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to count comments")
			return nil, err
		}

		likesCount, err := s.Repository.CountLikes(ctx, post.ID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to count likes")
			return nil, err
		}

		pl := &entity.PostList{
			PostID:   post.ID,
			UserID:   post.UserID,
			Title:    post.Title,
			Comments: commtentsCount,
			Likes:    likesCount,
		}

		postList = append(postList, pl)
	}

	return postList, nil
}

func (s *Service) SearchPosts(ctx context.Context, userID int, searchQuery SearchQuery) (*PaginatedResponse, error) {
	logger := s.Logger.With().Str("method", "SearchPosts").Str("query", searchQuery.Query).Logger()

	if err := searchQuery.Validate(); err != nil {
		logger.Error().Err(err).Msg("failed to validate search query")
		return nil, err
	}

	searchQuery.SetDefaults()

	posts, total, err := s.Repository.SearchPosts(ctx, userID, searchQuery.Query, searchQuery.GetOffset(), searchQuery.GetLimit())
	if err != nil {
		logger.Error().Err(err).Msg("failed to search posts")
		return nil, err
	}

	pagination := PaginationQuery{
		Page:     searchQuery.Page,
		PageSize: searchQuery.PageSize,
	}

	response := NewPaginatedResponse(posts, pagination, total)
	return response, nil
}

func (s *Service) GetPostByID(ctx context.Context, userID, postID int) (*entity.Post, error) {
	logger := s.Logger.With().Str("method", "GetPostByID").Int("user_id", userID).Int("post_id", postID).Logger()
	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	if postID <= 0 {
		logger.Error().Msg("post id is required")
		return nil, ErrPostIdRequired
	}

	post, err := s.Repository.GetPostByID(ctx, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get post by id")
		return nil, err
	}

	if post.UserID != userID && !post.IsPublic {
		return nil, ErrUnauthorized
	}

	comments, err := s.Repository.GetCommentsByCommentLikesDesc(ctx, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get comments by comment likes desc")
		return nil, err
	}

	likes, err := s.Repository.GetLikesByPostID(ctx, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get likes by post id")
		return nil, err
	}

	post.Comments = append(post.Comments, comments...)
	post.Likes = append(post.Likes, likes...)

	return post, nil
}

func (s *Service) UpdatePost(ctx context.Context, in UpdatePostIn, userID, postID int) (*entity.Post, error) {
	logger := s.Logger.With().Str("method", "UpdatePost").Int("user_id", userID).Int("post_id", postID).Logger()
	if err := in.Validate(); err != nil {
		logger.Error().Err(err).Msg("failed to validate update post")
		return nil, err
	}

	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return nil, ErrUserIdRequired
	}

	if postID <= 0 {
		logger.Error().Msg("post id is required")
		return nil, ErrPostIdRequired
	}

	// Проверяем, что пост принадлежит пользователю
	existingPost, err := s.Repository.GetPostByID(ctx, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get existing post")
		return nil, err
	}

	if existingPost.UserID != userID {
		logger.Error().Msg("user can only update own posts")
		return nil, ErrUnauthorized
	}

	updatedPost, err := s.Repository.UpdatePost(ctx, entity.Post{
		ID:       postID,
		UserID:   userID,
		Title:    in.Title,
		ImageURL: in.ImageURL,
		Content:  in.Content,
	})
	if err != nil {
		logger.Error().Err(err).Msg("failed to update post")
		return nil, err
	}

	return updatedPost, nil
}

func (s *Service) DeletePost(ctx context.Context, userID, postID int) error {
	logger := s.Logger.With().Str("method", "DeletePost").Int("user_id", userID).Int("post_id", postID).Logger()
	if userID <= 0 {
		logger.Error().Msg("user id is required")
		return ErrUserIdRequired
	}

	if postID <= 0 {
		return ErrPostIdRequired
	}

	// Проверяем, что пост принадлежит пользователю
	existingPost, err := s.Repository.GetPostByID(ctx, postID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to get existing post")
		return err
	}

	if existingPost.UserID != userID {
		logger.Error().Msg("user can only delete own posts")
		return ErrUnauthorized
	}

	tNow := time.Now()
	toUpdate := entity.Post{
		ID:        postID,
		DeletedAt: &tNow,
	}

	err = s.Repository.DeletePost(ctx, toUpdate)
	if err != nil {
		logger.Error().Err(err).Msg("failed to delete post")
		return err
	}

	return nil
}
