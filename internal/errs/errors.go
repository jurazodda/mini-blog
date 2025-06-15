package errs

import "errors"

var (
	// common
	ErrRecordNotFound = errors.New("ErrRecordNotFound")

	// users
	ErrUserNotFound = errors.New("ErrUserNotFound")
)