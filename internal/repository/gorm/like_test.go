package gorm

import (
	"context"
	"regexp"
	"testing"
	"time"

	"mini-blog/entity"
	"mini-blog/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupLikeTestDB(t *testing.T) (*GormRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}
	logg, _ := logger.NewLogger()
	repo := &GormRepository{DB: gormDB, Logger: logg}
	return repo, mock
}

func TestGetLikesByPostID_Success(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "likes" WHERE post_id = $1 AND deleted_at IS NULL`)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "created_at", "deleted_at"}).AddRow(1, 1, 1, now, nil))
	ctx := context.Background()
	likes, err := repo.GetLikesByPostID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(likes) != 1 {
		t.Fatalf("expected 1 like, got %d", len(likes))
	}
}

func TestGetLikesByPostID_Error(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "likes" WHERE post_id = $1 AND deleted_at IS NULL`)).WithArgs(1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetLikesByPostID(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCountLikes_Success(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "likes" WHERE post_id = $1 AND deleted_at IS NULL`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	ctx := context.Background()
	count, err := repo.CountLikes(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestCountLikes_Error(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "likes" WHERE post_id = $1 AND deleted_at IS NULL`)).WithArgs(1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.CountLikes(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLikeByUserAndPost_Success(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "likes" WHERE user_id = $1 AND post_id = $2 AND deleted_at IS NULL ORDER BY "likes"."id" LIMIT $3`)).WithArgs(1, 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "created_at", "deleted_at"}).AddRow(1, 1, 1, now, nil))
	ctx := context.Background()
	like, err := repo.GetLikeByUserAndPost(ctx, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if like.ID != 1 {
		t.Fatalf("expected like ID 1, got %d", like.ID)
	}
}

func TestGetLikeByUserAndPost_Error(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "likes" WHERE user_id = $1 AND post_id = $2 AND deleted_at IS NULL ORDER BY "likes"."id" LIMIT $3`)).WithArgs(1, 1, 1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetLikeByUserAndPost(ctx, 1, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateLike_Success(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	like := entity.Like{ID: 1, UserID: 1, PostID: 1}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "likes"`)).
		WithArgs(like.UserID, like.PostID, sqlmock.AnyArg(), nil, like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.CreateLike(ctx, like)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateLike_Error(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	like := entity.Like{ID: 1, UserID: 1, PostID: 1}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "likes"`)).
		WithArgs(like.UserID, like.PostID, sqlmock.AnyArg(), nil, like.ID).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.CreateLike(ctx, like)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteLike_Success(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	like := entity.Like{ID: 1, UserID: 1, PostID: 1}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "likes" WHERE`)).WithArgs(like.UserID, like.PostID, like.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	ctx := context.Background()
	err := repo.DeleteLike(ctx, like)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteLike_Error(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	like := entity.Like{ID: 1, UserID: 1, PostID: 1}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "likes" WHERE`)).WithArgs(like.UserID, like.PostID).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	err := repo.DeleteLike(ctx, like)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetLikeByUserAndComment_Success(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comment_likes" WHERE user_id = $1 AND comment_id = $2 ORDER BY "comment_likes"."id" LIMIT $3`)).WithArgs(1, 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "comment_id", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, 1, now, now, nil))
	ctx := context.Background()
	like, err := repo.GetLikeByUserAndComment(ctx, 1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if like.ID != 1 {
		t.Fatalf("expected like ID 1, got %d", like.ID)
	}
}

func TestGetLikeByUserAndComment_Error(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comment_likes" WHERE user_id = $1 AND comment_id = $2`)).WithArgs(1, 1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetLikeByUserAndComment(ctx, 1, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreateLikeComment_Success(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	like := entity.CommentLike{ID: 1, UserID: 1, CommentID: 1}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comment_likes"`)).WithArgs(like.UserID, like.CommentID, sqlmock.AnyArg(), nil, like.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.CreateLikeComment(ctx, like)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateLikeComment_Error(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	like := entity.CommentLike{ID: 1, UserID: 1, CommentID: 1}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comment_likes"`)).WithArgs(like.UserID, like.CommentID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.CreateLikeComment(ctx, like)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteLikeComment_Success(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "comment_likes"`)).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	ctx := context.Background()
	err := repo.DeleteLikeComment(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteLikeComment_Error(t *testing.T) {
	repo, mock := setupLikeTestDB(t)
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "comment_likes"`)).WithArgs(1).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	err := repo.DeleteLikeComment(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
