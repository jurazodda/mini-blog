package gorm

import (
	"errors"
	"mini-blog/internal/errs"

	"gorm.io/gorm"
)

func translateError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errs.ErrRecordNotFound
	}
	return nil
}