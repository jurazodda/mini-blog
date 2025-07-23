package controller

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"mini-blog/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestHandleError(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		expectCode int
		expectBody string
	}{
		{"record not found", gorm.ErrRecordNotFound, http.StatusNotFound, "Record not found"},
		{"user id required", service.ErrUserIdRequired, http.StatusBadRequest, "User ID is required"},
		{"post id required", service.ErrPostIdRequired, http.StatusBadRequest, "Post ID is required"},
		{"comment id required", service.ErrCommentIdRequired, http.StatusBadRequest, "Comment ID is required"},
		{"title required", service.ErrTitleRequired, http.StatusBadRequest, "Title is required"},
		{"content required", service.ErrContentRequired, http.StatusBadRequest, "Content is required"},
		{"comment required", service.ErrCommentRequired, http.StatusBadRequest, "Comment is required"},
		{"unauthorized", service.ErrUnauthorized, http.StatusUnauthorized, "Unauthorized access"},
		{"forbidden", service.ErrForbidden, http.StatusForbidden, "Forbidden access"},
		{"access denied", service.ErrAccessDenied, http.StatusForbidden, "Access denied"},
		{"username required", service.ErrUsernameRequired, http.StatusBadRequest, "Username is required"},
		{"email required", service.ErrEmailRequired, http.StatusBadRequest, "Email is required"},
		{"password required", service.ErrPasswordRequired, http.StatusBadRequest, "Password is required"},
		{"invalid email", service.ErrInvalidEmailFormat, http.StatusBadRequest, "Invalid email format"},
		{"invalid password", service.ErrInvalidPassword, http.StatusUnauthorized, "Invalid password"},
		{"unknown error", errors.New("some error"), http.StatusInternalServerError, "Internal server error"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			HandleError(c, tc.err)
			assert.Equal(t, tc.expectCode, w.Code)
			assert.Contains(t, w.Body.String(), tc.expectBody)
		})
	}
}
