package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchQuery_Validate(t *testing.T) {
	q := SearchQuery{Query: ""}
	assert.Error(t, q.Validate())
	q.Query = "test"
	assert.NoError(t, q.Validate())
}

func TestSearchQuery_SetDefaults(t *testing.T) {
	q := SearchQuery{}
	q.SetDefaults()
	assert.Equal(t, 1, q.Page)
	assert.Equal(t, 10, q.PageSize)
	q.Page = -5
	q.PageSize = 200
	q.SetDefaults()
	assert.Equal(t, 1, q.Page)
	assert.Equal(t, 100, q.PageSize)
}

func TestSearchQuery_GetOffsetAndLimit(t *testing.T) {
	q := SearchQuery{Page: 2, PageSize: 10}
	assert.Equal(t, 10, q.GetOffset())
	assert.Equal(t, 10, q.GetLimit())
}

func TestCreatePostIn_Validate(t *testing.T) {
	in := CreatePostIn{}
	assert.Error(t, in.Validate())
	in.Title = "t"
	in.Content = "c"
	assert.NoError(t, in.Validate())
}

func TestUpdatePostIn_Validate(t *testing.T) {
	in := UpdatePostIn{}
	assert.Error(t, in.Validate())
	in.Title = "t"
	in.Content = "c"
	assert.NoError(t, in.Validate())
}

func TestCreateCommentIn_Validate(t *testing.T) {
	in := CreateCommentIn{}
	assert.Error(t, in.Validate())
	in.PostID = 1
	in.Comment = "c"
	assert.NoError(t, in.Validate())
}

func TestCreateRepostIn_Validate(t *testing.T) {
	in := CreateRepostIn{}
	assert.Error(t, in.Validate())
	in.PostID = 1
	assert.NoError(t, in.Validate())
}

func TestUpdateRepostIn_Validate(t *testing.T) {
	in := UpdateRepostIn{}
	assert.Error(t, in.Validate())
	in.Title = "t"
	in.Content = "c"
	assert.NoError(t, in.Validate())
}

func TestSignUpIn_Validate(t *testing.T) {
	in := SignUpIn{}
	assert.Error(t, in.Validate())
	in.Username = "u"
	in.Email = "bad"
	in.Password = "123"
	assert.Error(t, in.Validate())
	in.Email = "user@example.com"
	in.Password = "password123"
	assert.NoError(t, in.Validate())
}

func TestSignInIn_Validate(t *testing.T) {
	in := SignInIn{}
	assert.Error(t, in.Validate())
	in.Email = "bad"
	in.Password = "123"
	assert.Error(t, in.Validate())
	in.Email = "user@example.com"
	in.Password = "password123"
	assert.NoError(t, in.Validate())
}
