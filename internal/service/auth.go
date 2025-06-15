package service

import (
	"context"
	"errors"
	"mini-blog/internal/entity"
	"mini-blog/internal/errs"
	"mini-blog/pkg/security"
	"regexp"
	"unicode/utf8"
)

const (
	/*
		Email
		- Имя пользователя перед @
		- Наличие @ — символ "собачки"
		- Доменное имя
		- Наличие точки
		— Доменная зона, минимум 2 символа
	*/
	regexpEmail = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)

type SignUpIn struct {
	Username string
	Email    string
	Password string
}

func (in *SignUpIn) Validate() error {
	if in.Username == "" {
		return errors.New("username is required")
	}

	if in.Email == "" {
		return errors.New("email is required")
	}

	emailRegex := regexp.MustCompile(regexpEmail)
	if !emailRegex.MatchString(in.Email) {
		return errors.New("invalid email format")
	}

	if in.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (s *Service) SignUp(ctx context.Context, in SignUpIn) error {
	err := in.Validate()
	if err != nil {
		return errors.Join(errors.New("ErrValidationFailed"), err)
	}

	passHash, err := security.GeneratePasswordHash(in.Password)
	if err != nil {
		return err
	}

	err = s.Repository.CreateUser(ctx, entity.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: passHash,
	})
	if err != nil {
		return err
	}

	return nil
}

type SignInIn struct {
	Email    string
	Password string
}

func (in *SignInIn) Validate() error {
	if in.Email == "" {
		return errors.New("username is required")
	}

	emailRegex := regexp.MustCompile(regexpEmail)
	if !emailRegex.MatchString(in.Email) {
		return errors.New("invalid email format")
	}

	if in.Password == "" {
		return errors.New("password is required")
	}

	length := utf8.RuneCountInString(in.Password)
	if length < 8 || length > 16 {
		return errors.New("password must be between 8 and 16 characters")
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, in SignInIn) (*entity.User, error) {
	logger := s.Logger.With().Ctx(ctx).Str("method", "SignIn").Str("email", in.Email).Logger()

	err := in.Validate()
	if err != nil {
		logger.Error().Err(err).Msg("validation failed")
		return nil, errors.Join(errors.New("ErrValidationFailed"), err)
	}

	user, err := s.Repository.GetUserByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return nil, errs.ErrUserNotFound
		}
		logger.Error().Err(err).Msg("failed getting user by email")
		return nil, err
	}

	if !security.ComparePassword(user.PasswordHash, in.Password) {
		logger.Error().Err(err).Msg("failed compare password")
		return nil, errors.New("invalid password")
	}

	return user, nil
}

