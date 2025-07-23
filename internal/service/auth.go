package service

import (
	"context"
	"errors"
	"mini-blog/entity"
	"mini-blog/pkg/password"
)

func (s *Service) CreateUser(ctx context.Context, in CreateUserIn) (*entity.User, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	passHash, err := password.GeneratePasswordHash(in.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.Repository.CreateUser(ctx, entity.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: passHash,
	})
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, in GetUserIn) (*entity.User, error) {
	err := in.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.Repository.GetUserByID(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) LoginUser(ctx context.Context, in LoginUserIn) (*LoginResponse, error) {
	log := s.Logger.Info().Str("method", "LoginUser").Str("email", in.Email)

	err := in.Validate()
	if err != nil {
		s.Logger.Error().Err(err).Msg("validation failed")
		return nil, err
	}

	user, err := s.Repository.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		s.Logger.Error().Err(err).Msg("failed getting user by email")
		return nil, err
	}

	if !password.ComparePassword(user.PasswordHash, in.Password) {
		s.Logger.Error().Err(err).Msg("failed compare password")
		return nil, ErrInvalidPassword
	}

	// Generate JWT token
	token, err := GenerateToken(ctx, user.ID)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to generate token")
		return nil, err
	}

	log.Msg("user logged in")
	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *Service) SignUp(ctx context.Context, in SignUpIn) error {
	err := in.Validate()
	if err != nil {
		return err
	}

	passHash, err := password.GeneratePasswordHash(in.Password)
	if err != nil {
		return err
	}

	_, err = s.Repository.CreateUser(ctx, entity.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: passHash,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, in SignInIn) (*entity.User, error) {
	log := s.Logger.Info().Str("method", "SignIn").Str("email", in.Email)

	err := in.Validate()
	if err != nil {
		s.Logger.Error().Err(err).Msg("validation failed")
		return nil, err
	}

	user, err := s.Repository.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		s.Logger.Error().Err(err).Msg("failed getting user by email")
		return nil, err
	}

	if !password.ComparePassword(user.PasswordHash, in.Password) {
		s.Logger.Error().Err(err).Msg("failed compare password")
		return nil, ErrInvalidPassword
	}

	log.Msg("user signed in")
	return user, nil
}
