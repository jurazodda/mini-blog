package service

import (
	"context"
	"mini-blog/entity"
	"mini-blog/internal/repository"
	"mini-blog/pkg/logger"
)

// ServiceInterface defines the contract for business logic layer
type ServiceI interface {
	// Auth methods (legacy)
	SignUp(ctx context.Context, in SignUpIn) error
	SignIn(ctx context.Context, in SignInIn) (*entity.User, error)

	// User methods
	CreateUser(ctx context.Context, in CreateUserIn) (*entity.User, error)
	GetUser(ctx context.Context, in GetUserIn) (*entity.User, error)
	LoginUser(ctx context.Context, in LoginUserIn) (*LoginResponse, error)

	// Post methods
	CreatePost(ctx context.Context, in CreatePostIn, userID int) (*entity.Post, error)
	GetPosts(ctx context.Context, userID int) ([]*entity.PostList, error)
	SearchPosts(ctx context.Context, userID int, searchQuery SearchQuery) (*PaginatedResponse, error)
	GetPostByID(ctx context.Context, userID, postID int) (*entity.Post, error)
	UpdatePost(ctx context.Context, in UpdatePostIn, userID, postID int) (*entity.Post, error)
	DeletePost(ctx context.Context, userID, postID int) error

	// Comment methods
	CreateComment(ctx context.Context, in CreateCommentIn, userID int) (*entity.Comment, error)
	GetCommentByID(ctx context.Context, userID, commentID int) (*entity.Comment, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]*entity.Comment, error)
	UpdateComment(ctx context.Context, comment string, userID, commentID int) (*entity.Comment, error)
	DeleteComment(ctx context.Context, userID, commentID int) error

	// Post like methods
	CreateLikePost(ctx context.Context, postID, userID int) (*ToggleLikeResult, error)
	CreateLikeComment(ctx context.Context, commentID, userID int) (*ToggleLikeResult, error)

	// Repost methods
	CreateRepost(ctx context.Context, in CreateRepostIn, userID int) (*entity.Repost, error)
	GetReposts(ctx context.Context) ([]*entity.Repost, error)
	GetRepostByID(ctx context.Context, userID, repostID int) (*entity.Repost, error)
	GetRepostsByUserID(ctx context.Context, targetUserID int) ([]*entity.Repost, error)
	UpdateRepost(ctx context.Context, in UpdateRepostIn, userID, repostID int) (*entity.Repost, error)
	DeleteRepost(ctx context.Context, userID, repostID int) error
}

type Service struct {
	Repository repository.RepositoryI
	Logger     *logger.Logger
}

// Ensure Service implements ServiceInterface
var _ ServiceI = (*Service)(nil)

func NewService(repository repository.RepositoryI, logger *logger.Logger) *Service {
	return &Service{
		Repository: repository,
		Logger:     logger,
	}
}
