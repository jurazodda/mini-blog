package repository

import (
	"context"
	"mini-blog/entity"
)

// Repository defines the contract for data access layer
type RepositoryI interface {
	// User methods
	CreateUser(ctx context.Context, user entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, userID int) (*entity.User, error)

	// Post methods
	CreatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	GetPosts(ctx context.Context, postID int) ([]*entity.Post, error)
	SearchPosts(ctx context.Context, userID int, query string, offset, limit int) ([]*entity.Post, int64, error)
	GetPostByID(ctx context.Context, postID int) (*entity.Post, error)
	UpdatePost(ctx context.Context, post entity.Post) (*entity.Post, error)
	DeletePost(ctx context.Context, post entity.Post) error

	// Comment methods
	CreateComment(ctx context.Context, comment entity.Comment) (*entity.Comment, error)
	GetCommentByID(ctx context.Context, commentID int) (*entity.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]*entity.Comment, error)
	GetCommentsByCommentLikesDesc(ctx context.Context, postID int) ([]*entity.Comment, error)
	UpdateComment(ctx context.Context, comment entity.Comment) (*entity.Comment, error)
	DeleteComment(ctx context.Context, comment entity.Comment) error
	CountComments(ctx context.Context, postID int) (int, error)

	// Post like methods
	CreateLike(ctx context.Context, like entity.Like) (*entity.Like, error)
	GetLikesByPostID(ctx context.Context, postID int) ([]*entity.Like, error)
	GetLikeByUserAndPost(ctx context.Context, userID, postID int) (*entity.Like, error)
	DeleteLike(ctx context.Context, like entity.Like) error
	CountLikes(ctx context.Context, postID int) (int, error)

	// Comment like methods
	CreateLikeComment(ctx context.Context, like entity.CommentLike) (*entity.CommentLike, error)
	GetLikeByUserAndComment(ctx context.Context, userID, commentID int) (*entity.CommentLike, error)
	DeleteLikeComment(ctx context.Context, likeID int) error

	// Repost methods
	CreateRepost(ctx context.Context, repost entity.Repost) (*entity.Repost, error)
	GetReposts(ctx context.Context) ([]*entity.Repost, error)
	GetRepostByID(ctx context.Context, repostID int) (*entity.Repost, error)
	GetRepostsByUserID(ctx context.Context, userID int) ([]*entity.Repost, error)
	UpdateRepost(ctx context.Context, repost entity.Repost) (*entity.Repost, error)
	DeleteRepost(ctx context.Context, repost entity.Repost) error
}
