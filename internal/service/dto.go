package service

import (
	"errors"
	"regexp"
	"unicode/utf8"
)

type CreateCommentIn struct {
	PostID  int    `json:"post_id"`
	Comment string `json:"comment"`
}

func (in *CreateCommentIn) Validate() error {
	if in.PostID <= 0 {
		return ErrPostIdRequired
	}

	if in.Comment == "" {
		return ErrCommentRequired
	}

	return nil
}

type CreatePostIn struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	Content  string `json:"content"`
}

func (in *CreatePostIn) Validate() error {
	if in.Title == "" {
		return ErrTitleRequired
	}

	// validate imaeUrl

	if in.Content == "" {
		return ErrContentRequired
	}

	return nil
}

type UpdatePostIn struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	Content  string `json:"content"`
}

func (in *UpdatePostIn) Validate() error {
	if in.Title == "" {
		return ErrTitleRequired
	}

	// validate imaeUrl

	if in.Content == "" {
		return ErrContentRequired
	}

	return nil
}

type CreateRepostIn struct {
	PostID  int    `json:"post_id"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
}

func (in *CreateRepostIn) Validate() error {
	if in.PostID <= 0 {
		return ErrPostIdRequired
	}
	return nil
}

type UpdateRepostIn struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (in *UpdateRepostIn) Validate() error {
	if in.Title == "" {
		return ErrTitleRequired
	}
	if in.Content == "" {
		return ErrContentRequired
	}
	return nil
}

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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (in *SignUpIn) Validate() error {
	if in.Username == "" {
		return ErrUsernameRequired
	}

	if in.Email == "" {
		return ErrEmailRequired
	}

	emailRegex := regexp.MustCompile(regexpEmail)
	if !emailRegex.MatchString(in.Email) {
		return ErrInvalidEmailFormat
	}

	if in.Password == "" {
		return ErrPasswordRequired
	}

	length := utf8.RuneCountInString(in.Password)
	if length < 8 || length > 16 {
		return ErrInvalidPassword
	}

	return nil
}

type SignInIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (in *SignInIn) Validate() error {
	if in.Email == "" {
		return ErrEmailRequired
	}

	emailRegex := regexp.MustCompile(regexpEmail)
	if !emailRegex.MatchString(in.Email) {
		return ErrInvalidEmailFormat
	}

	if in.Password == "" {
		return ErrPasswordRequired
	}

	length := utf8.RuneCountInString(in.Password)
	if length < 8 || length > 16 {
		return errors.New("password must be between 8 and 16 characters")
	}

	return nil
}

// Additional DTOs for user operations
type CreateUserIn = SignUpIn
type GetUserIn struct {
	UserID int `json:"user_id"`
}

func (in *GetUserIn) Validate() error {
	if in.UserID <= 0 {
		return ErrUserIdRequired
	}
	return nil
}

type LoginUserIn = SignInIn

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
