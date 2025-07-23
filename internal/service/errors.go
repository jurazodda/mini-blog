package service

import "errors"

var (
	// auth
	ErrUsernameRequired   = errors.New("username is required")
	ErrEmailRequired      = errors.New("email is required")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrPasswordRequired   = errors.New("password is required")
	ErrInvalidPassword    = errors.New("invalid password")

	// authorization & security
	ErrUnauthorized = errors.New("unauthorized access")
	ErrForbidden    = errors.New("forbidden access")
	ErrAccessDenied = errors.New("access denied")

	// post
	ErrPostIdRequired  = errors.New("post id required")
	ErrUserIdRequired  = errors.New("user id required")
	ErrTitleRequired   = errors.New("title is required")
	ErrContentRequired = errors.New("content is required")

	// comment
	ErrCommentRequired   = errors.New("comment is required")
	ErrCommentIdRequired = errors.New("comment id required")

	// repost
	ErrRepostIdRequired = errors.New("repost id required")

	// common
	ErrRecordNotFound = errors.New("ErrRecordNotFound")

	// users
	ErrUserNotFound = errors.New("ErrUserNotFound")
)
