package service

import (
	"context"
	"errors"
	"mini-blog/entity"
	"mini-blog/internal/repository"
	"mini-blog/mocks"
	"mini-blog/pkg/logger"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type PostServiceTestSuite struct {
	suite.Suite
	service        ServiceI
	repository     repository.RepositoryI
	mockRepository *mocks.MockRepositoryI
	ctx            context.Context
}

func (s *PostServiceTestSuite) SetupTest() {
	s.mockRepository = mocks.NewMockRepositoryI(s.T())
	s.repository = s.mockRepository
	log, err := logger.NewLogger()
	s.Require().NoError(err)
	s.service = NewService(s.repository, log)
	s.ctx = context.Background()
	t := s.T()
	t.Setenv("SIGNING_KEY", "test-signing-key-for-tests")
	t.Setenv("TOKEN_TTL", "24")
}

func (s *PostServiceTestSuite) TearDownTest() {
	os.Unsetenv("SIGNING_KEY")
	os.Unsetenv("TOKEN_TTL")
}

func TestPostServiceTestSuite(t *testing.T) {
	suite.Run(t, new(PostServiceTestSuite))
}

func (s *PostServiceTestSuite) createTestPost() *entity.Post {
	return &entity.Post{
		ID:        1,
		Title:     "Test Post",
		Content:   "Test content",
		ImageURL:  "http://example.com/image.jpg",
		UserID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *PostServiceTestSuite) createTestComment() *entity.Comment {
	return &entity.Comment{
		ID:        1,
		PostID:    1,
		UserID:    1,
		Comment:   "Test comment",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *PostServiceTestSuite) TestCreatePost_Success() {
	input := CreatePostIn{
		Title:    "New Post",
		Content:  "Post content",
		ImageURL: "http://example.com/image.jpg",
	}
	userID := 1
	expectedPost := &entity.Post{
		ID:       1,
		Title:    input.Title,
		Content:  input.Content,
		ImageURL: input.ImageURL,
		UserID:   userID,
	}
	s.mockRepository.EXPECT().CreatePost(mock.Anything, entity.Post{
		Title:    input.Title,
		Content:  input.Content,
		ImageURL: input.ImageURL,
		UserID:   userID,
	}).Return(expectedPost, nil).Once()
	result, err := s.service.CreatePost(s.ctx, input, userID)
	s.NoError(err)
	s.Equal(expectedPost, result)
}

func (s *PostServiceTestSuite) TestCreatePost_ValidationError() {
	input := CreatePostIn{
		Title:   "",
		Content: "",
	}
	userID := 1
	result, err := s.service.CreatePost(s.ctx, input, userID)
	s.Error(err)
	s.Nil(result)
	s.Contains(err.Error(), "title is required")
}

func (s *PostServiceTestSuite) TestGetPosts_Success() {
	posts := []*entity.Post{s.createTestPost()}
	s.mockRepository.EXPECT().GetPosts(mock.Anything, mock.Anything).Return(posts, nil).Once()
	s.mockRepository.EXPECT().CountComments(mock.Anything, 1).Return(5, nil).Once()
	s.mockRepository.EXPECT().CountLikes(mock.Anything, 1).Return(10, nil).Once()
	result, err := s.service.GetPosts(s.ctx, 1)
	s.NoError(err)
	s.Len(result, 1)
	s.Equal(5, result[0].Comments)
	s.Equal(10, result[0].Likes)
}

func (s *PostServiceTestSuite) TestGetPostByID_Success() {
	userID, postID := 1, 1
	post := s.createTestPost()
	comments := []*entity.Comment{s.createTestComment()}
	likes := []*entity.Like{{ID: 1, UserID: 1, PostID: 1}}
	s.mockRepository.EXPECT().GetPostByID(mock.Anything, postID).Return(post, nil).Once()
	s.mockRepository.EXPECT().GetCommentsByCommentLikesDesc(mock.Anything, postID).Return(comments, nil).Once()
	s.mockRepository.EXPECT().GetLikesByPostID(mock.Anything, postID).Return(likes, nil).Once()
	result, err := s.service.GetPostByID(s.ctx, userID, postID)
	s.NoError(err)
	s.Equal(post.ID, result.ID)
	s.Len(result.Comments, 1)
	s.Len(result.Likes, 1)
}

func (s *PostServiceTestSuite) TestUpdatePost_Success() {
	input := UpdatePostIn{
		Title:    "Updated Title",
		Content:  "Updated content",
		ImageURL: "http://example.com/updated.jpg",
	}
	userID, postID := 1, 1
	existingPost := s.createTestPost()
	existingPost.UserID = userID
	updatedPost := &entity.Post{
		ID:       postID,
		Title:    input.Title,
		Content:  input.Content,
		ImageURL: input.ImageURL,
		UserID:   userID,
	}
	s.mockRepository.EXPECT().GetPostByID(mock.Anything, postID).Return(existingPost, nil).Once()
	s.mockRepository.EXPECT().UpdatePost(mock.Anything, entity.Post{
		ID:       postID,
		Title:    input.Title,
		Content:  input.Content,
		ImageURL: input.ImageURL,
		UserID:   userID,
	}).Return(updatedPost, nil).Once()
	result, err := s.service.UpdatePost(s.ctx, input, userID, postID)
	s.NoError(err)
	s.Equal(updatedPost, result)
}

func (s *PostServiceTestSuite) TestUpdatePost_Unauthorized() {
	input := UpdatePostIn{
		Title:   "Updated Title",
		Content: "Updated content",
	}
	userID, postID := 1, 1
	existingPost := s.createTestPost()
	existingPost.UserID = 2
	s.mockRepository.EXPECT().GetPostByID(mock.Anything, postID).Return(existingPost, nil).Once()
	result, err := s.service.UpdatePost(s.ctx, input, userID, postID)
	s.Error(err)
	s.Nil(result)
	s.True(errors.Is(err, ErrUnauthorized))
}

func (s *PostServiceTestSuite) TestDeletePost_Success() {
	userID, postID := 1, 1
	existingPost := s.createTestPost()
	existingPost.UserID = userID
	s.mockRepository.EXPECT().GetPostByID(mock.Anything, postID).Return(existingPost, nil).Once()
	s.mockRepository.EXPECT().DeletePost(mock.Anything, mock.MatchedBy(func(post entity.Post) bool {
		return post.ID == postID && post.DeletedAt != nil
	})).Return(nil).Once()
	err := s.service.DeletePost(s.ctx, userID, postID)
	s.NoError(err)
}
