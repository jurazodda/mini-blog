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

type AuthServiceTestSuite struct {
	suite.Suite
	service        ServiceI
	repository     repository.RepositoryI
	mockRepository *mocks.MockRepositoryI
	ctx            context.Context
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.mockRepository = mocks.NewMockRepositoryI(s.T())
	s.repository = s.mockRepository
	log, err := logger.NewLogger()
	s.Require().NoError(err)
	s.service = NewService(s.repository, log)
	s.ctx = context.Background()
	s.T().Setenv("SIGNING_KEY", "test-signing-key-for-tests")
	s.T().Setenv("TOKEN_TTL", "24")
}

func (s *AuthServiceTestSuite) TearDownTest() {
	os.Unsetenv("SIGNING_KEY")
	os.Unsetenv("TOKEN_TTL")
}

func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (s *AuthServiceTestSuite) createTestUser() *entity.User {
	return &entity.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "$2a$10$HEprfLqaBWP9rUsuQtq4.u0YZo7QanIXuOjLA3Wx95PwGeYZhobH2",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (s *AuthServiceTestSuite) TestCreateUser_Success() {
	input := CreateUserIn{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}
	expectedUser := &entity.User{
		ID:       1,
		Username: input.Username,
		Email:    input.Email,
	}
	s.mockRepository.EXPECT().CreateUser(mock.Anything, mock.MatchedBy(func(user entity.User) bool {
		return user.Username == input.Username &&
			user.Email == input.Email &&
			len(user.PasswordHash) > 0
	})).Return(expectedUser, nil).Once()
	result, err := s.service.CreateUser(s.ctx, input)
	s.NoError(err)
	s.Equal(expectedUser.Username, result.Username)
	s.Equal(expectedUser.Email, result.Email)
}

func (s *AuthServiceTestSuite) TestCreateUser_ValidationError() {
	input := CreateUserIn{
		Username: "",
		Email:    "invalid-email",
		Password: "",
	}
	result, err := s.service.CreateUser(s.ctx, input)
	s.Error(err)
	s.Nil(result)
}

func (s *AuthServiceTestSuite) TestLoginUser_UserNotFound() {
	input := LoginUserIn{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}
	s.mockRepository.EXPECT().GetUserByEmail(mock.Anything, input.Email).Return(nil, ErrRecordNotFound).Once()
	result, err := s.service.LoginUser(s.ctx, input)
	s.Error(err)
	s.Nil(result)
	s.True(errors.Is(err, ErrUserNotFound))
}

func (s *AuthServiceTestSuite) TestGetUser_Success() {
	input := GetUserIn{UserID: 1}
	user := s.createTestUser()
	s.mockRepository.EXPECT().GetUserByID(mock.Anything, input.UserID).Return(user, nil).Once()
	result, err := s.service.GetUser(s.ctx, input)
	s.NoError(err)
	s.Equal(user, result)
}

func (s *AuthServiceTestSuite) TestGetUser_ValidationError() {
	input := GetUserIn{UserID: 0}
	result, err := s.service.GetUser(s.ctx, input)
	s.Error(err)
	s.Nil(result)
}

func (s *AuthServiceTestSuite) TestSignUp_Success() {
	input := SignUpIn{
		Username: "newuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}
	s.mockRepository.EXPECT().CreateUser(mock.Anything, mock.AnythingOfType("entity.User")).Return(&entity.User{ID: 1}, nil).Once()
	err := s.service.SignUp(s.ctx, input)
	s.NoError(err)
}

func (s *AuthServiceTestSuite) TestSignUp_ValidationError() {
	input := SignUpIn{Username: "", Email: "bad", Password: ""}
	err := s.service.SignUp(s.ctx, input)
	s.Error(err)
}

func (s *AuthServiceTestSuite) TestSignUp_RepoError() {
	input := SignUpIn{Username: "user", Email: "user@example.com", Password: "password123"}
	s.mockRepository.EXPECT().CreateUser(mock.Anything, mock.AnythingOfType("entity.User")).Return(nil, errors.New("db error")).Once()
	err := s.service.SignUp(s.ctx, input)
	s.Error(err)
}

func (s *AuthServiceTestSuite) TestSignIn_Success() {
	input := SignInIn{Email: "user@example.com", Password: "password123"}
	user := s.createTestUser()
	s.mockRepository.EXPECT().GetUserByEmail(mock.Anything, input.Email).Return(user, nil).Once()
	result, err := s.service.SignIn(s.ctx, input)
	s.NoError(err)
	s.Equal(user, result)
}

func (s *AuthServiceTestSuite) TestSignIn_ValidationError() {
	input := SignInIn{Email: "", Password: ""}
	result, err := s.service.SignIn(s.ctx, input)
	s.Error(err)
	s.Nil(result)
}

func (s *AuthServiceTestSuite) TestSignIn_UserNotFound() {
	input := SignInIn{Email: "user@example.com", Password: "password123"}
	s.mockRepository.EXPECT().GetUserByEmail(mock.Anything, input.Email).Return(nil, ErrRecordNotFound).Once()
	result, err := s.service.SignIn(s.ctx, input)
	s.Error(err)
	s.Nil(result)
	s.True(errors.Is(err, ErrUserNotFound))
}

func (s *AuthServiceTestSuite) TestSignIn_InvalidPassword() {
	input := SignInIn{Email: "user@example.com", Password: "wrongpass"}
	user := s.createTestUser()
	s.mockRepository.EXPECT().GetUserByEmail(mock.Anything, input.Email).Return(user, nil).Once()
	result, err := s.service.SignIn(s.ctx, input)
	s.Error(err)
	s.Nil(result)
	s.True(errors.Is(err, ErrInvalidPassword))
}

func (s *AuthServiceTestSuite) TestSignIn_RepoError() {
	input := SignInIn{Email: "user@example.com", Password: "password123"}
	s.mockRepository.EXPECT().GetUserByEmail(mock.Anything, input.Email).Return(nil, errors.New("db error")).Once()
	result, err := s.service.SignIn(s.ctx, input)
	s.Error(err)
	s.Nil(result)
}
