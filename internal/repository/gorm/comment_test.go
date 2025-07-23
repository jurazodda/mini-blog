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

func setupCommentTestDB(t *testing.T) (*GormRepository, sqlmock.Sqlmock) {
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

func TestCreateComment_Success(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	comment := entity.Comment{ID: 1, PostID: 1, UserID: 1, Comment: "Test"}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comments"`)).
		WithArgs(comment.UserID, comment.PostID, comment.Comment, comment.CommentLike, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, comment.ID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.CreateComment(ctx, comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateComment_Error(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	comment := entity.Comment{ID: 1, PostID: 1, UserID: 1, Comment: "Test"}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "comments"`)).
		WithArgs(comment.UserID, comment.PostID, comment.Comment, comment.CommentLike, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, comment.ID).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.CreateComment(ctx, comment)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetCommentsByPostID_Success(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE post_id = $1 AND deleted_at IS NULL`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "comment_like", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, 1, "Test", 0, now, now, nil))
	ctx := context.Background()
	comments, err := repo.GetCommentsByPostID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
}

func TestGetCommentsByPostID_Error(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE post_id = $1 AND deleted_at IS NULL`)).WithArgs(1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetCommentsByPostID(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetCommentsByCommentLikesDesc_Success(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE post_id = $1 AND deleted_at IS NULL`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "comment_like", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, 1, "Test", 0, now, now, nil))
	ctx := context.Background()
	comments, err := repo.GetCommentsByCommentLikesDesc(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(comments))
	}
}

func TestGetCommentsByCommentLikesDesc_Error(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE post_id = $1 AND deleted_at IS NULL`)).WithArgs(1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetCommentsByCommentLikesDesc(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCountComments_Success(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "comments" WHERE post_id = $1 AND deleted_at IS NULL`)).WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	ctx := context.Background()
	count, err := repo.CountComments(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestCountComments_Error(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "comments" WHERE post_id = $1 AND deleted_at IS NULL`)).WithArgs(1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.CountComments(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetCommentByID_Success(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE deleted_at IS NULL AND "comments"."id" = $1 ORDER BY "comments"."id" LIMIT $2`)).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "comment_like", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, 1, "Test", 0, now, now, nil))
	ctx := context.Background()
	comment, err := repo.GetCommentByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != 1 {
		t.Fatalf("expected comment ID 1, got %d", comment.ID)
	}
}

func TestGetCommentByID_Error(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "comments" WHERE deleted_at IS NULL AND "comments"."id" = $1 ORDER BY "comments"."id" LIMIT $2`)).WithArgs(1, 1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetCommentByID(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdateComment_Success(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	comment := entity.Comment{
		ID: 1, PostID: 1, UserID: 1, Comment: "Updated",
	}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET`)).WithArgs(
		comment.ID, comment.UserID, comment.PostID, comment.Comment, sqlmock.AnyArg(), comment.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.UpdateComment(ctx, comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateComment_Error(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	comment := entity.Comment{
		ID: 1, PostID: 1, UserID: 1, Comment: "Updated",
	}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET`)).WithArgs(
		comment.ID, comment.UserID, comment.PostID, comment.Comment, sqlmock.AnyArg(), comment.ID,
	).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.UpdateComment(ctx, comment)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeleteComment_Success(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	comment := entity.Comment{ID: 1}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET`)).WithArgs(sqlmock.AnyArg(), comment.ID, comment.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	ctx := context.Background()
	err := repo.DeleteComment(ctx, comment)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteComment_Error(t *testing.T) {
	repo, mock := setupCommentTestDB(t)
	comment := entity.Comment{ID: 1}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "comments" SET`)).WithArgs(sqlmock.AnyArg(), comment.ID, comment.ID).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	err := repo.DeleteComment(ctx, comment)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
