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

func setupRepostTestDB(t *testing.T) (*GormRepository, sqlmock.Sqlmock) {
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

func TestCreateRepost_Success(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	repost := entity.Repost{ID: 1, UserID: 1, PostID: 1, Title: "Test"}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "reposts"`)).
		WithArgs(repost.UserID, repost.PostID, repost.Title, repost.Content, repost.ImageURL, repost.Comments, repost.Likes, sqlmock.AnyArg(), nil, repost.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.CreateRepost(ctx, repost)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateRepost_Error(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	repost := entity.Repost{ID: 1, UserID: 1, PostID: 1, Title: "Test"}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "reposts"`)).
		WithArgs(repost.UserID, repost.PostID, repost.Title, repost.Content, repost.ImageURL, repost.Comments, repost.Likes, sqlmock.AnyArg(), nil, repost.ID).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.CreateRepost(ctx, repost)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetReposts_Success(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reposts" WHERE deleted_at IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "title", "content", "image_url", "comments", "likes", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, 1, "Test", "", "", 0, 0, now, now, nil))
	ctx := context.Background()
	reposts, err := repo.GetReposts(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reposts) != 1 {
		t.Fatalf("expected 1 repost, got %d", len(reposts))
	}
}

func TestGetReposts_Error(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reposts" WHERE deleted_at IS NULL`)).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetReposts(ctx)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRepostByID_Success(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reposts" WHERE deleted_at IS NULL AND "reposts"."id" = $1 ORDER BY "reposts"."id" LIMIT $2`)).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "title", "content", "image_url", "comments", "likes", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, 1, "Test", "", "", 0, 0, now, now, nil))
	ctx := context.Background()
	repost, err := repo.GetRepostByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repost.ID != 1 {
		t.Fatalf("expected repost ID 1, got %d", repost.ID)
	}
}

func TestGetRepostByID_Error(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reposts" WHERE deleted_at IS NULL AND "reposts"."id" = $1 ORDER BY "reposts"."id" LIMIT $2`)).WithArgs(1, 1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetRepostByID(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func ptr(s string) *string { return &s }

func TestUpdateRepost_Success(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	now := time.Now()
	repost := entity.Repost{ID: 1, PostID: 1, UserID: 1, Title: "Updated", Content: "", ImageURL: nil, Comments: 0, Likes: 0, CreatedAt: now, DeletedAt: nil}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reposts" SET`)).WithArgs(
		repost.ID, repost.UserID, repost.PostID, repost.Title, repost.CreatedAt, repost.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.UpdateRepost(ctx, repost)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateRepost_Error(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	now := time.Now()
	repost := entity.Repost{ID: 1, PostID: 1, UserID: 1, Title: "Updated", Content: "", ImageURL: nil, Comments: 0, Likes: 0, CreatedAt: now, DeletedAt: nil}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reposts" SET`)).WithArgs(
		repost.ID, repost.UserID, repost.PostID, repost.Title, repost.CreatedAt, repost.ID,
	).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.UpdateRepost(ctx, repost)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteRepost_Success(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	repost := entity.Repost{ID: 1, PostID: 1, UserID: 1, Title: "ToDelete"}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reposts" SET`)).WithArgs(repost.UserID, repost.PostID, repost.Title, repost.ID, repost.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	ctx := context.Background()
	err := repo.DeleteRepost(ctx, repost)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteRepost_Error(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	repost := entity.Repost{ID: 1, PostID: 1, UserID: 1, Title: "ToDelete"}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reposts" SET`)).WithArgs(repost.UserID, repost.PostID, repost.Title, repost.ID, repost.ID).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	err := repo.DeleteRepost(ctx, repost)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRepostsByUserID_Success(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reposts" WHERE user_id = $1 AND deleted_at IS NULL`)).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at", "image_url"}).AddRow(1, 1, 1, "Test", "Content", now, now, nil, "url"))
	ctx := context.Background()
	reposts, err := repo.GetRepostsByUserID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(reposts) != 1 {
		t.Fatalf("expected 1 repost, got %d", len(reposts))
	}
}

func TestGetRepostsByUserID_Error(t *testing.T) {
	repo, mock := setupRepostTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reposts" WHERE user_id = $1 AND deleted_at IS NULL`)).WithArgs(1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetRepostsByUserID(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
