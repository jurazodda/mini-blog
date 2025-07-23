package gorm

import (
	"context"
	"mini-blog/entity"
	"regexp"
	"testing"
	"time"

	"mini-blog/pkg/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*GormRepository, sqlmock.Sqlmock) {
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

func TestCreatePost_Success(t *testing.T) {
	repo, mock := setupTestDB(t)
	post := entity.Post{ID: 1, Title: "Test", Content: "Content", UserID: 1}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "posts"`)).
		WithArgs(
			post.UserID,      // user_id
			post.Title,       // title
			post.ImageURL,    // image_url
			post.Content,     // content
			post.IsPublic,    // is_public
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			nil,              // deleted_at
			post.ID,          // id
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.CreatePost(ctx, post)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreatePost_Error(t *testing.T) {
	repo, mock := setupTestDB(t)
	post := entity.Post{ID: 1, Title: "Test", Content: "Content", UserID: 1}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "posts"`)).
		WithArgs(post.UserID, post.Title, post.ImageURL, post.Content, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, post.ID).
		WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.CreatePost(ctx, post)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPosts_Success(t *testing.T) {
	userID := 1
	repo, mock := setupTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE deleted_at IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "image_url", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 1, "Test", "Content", "url", now, now, nil))
	ctx := context.Background()
	posts, err := repo.GetPosts(ctx, userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
}

func TestGetPosts_Error(t *testing.T) {
	userID := 1
	ctx := context.Background()
	repo, mock := setupTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE deleted_at IS NULL`)).WillReturnError(gorm.ErrInvalidDB)
	_, err := repo.GetPosts(ctx, userID)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPostByID_Success(t *testing.T) {
	repo, mock := setupTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE deleted_at IS NULL AND "posts"."id" = $1 ORDER BY "posts"."id" LIMIT $2`)).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "image_url", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 1, "Test", "Content", "url", now, now, nil))
	ctx := context.Background()
	post, err := repo.GetPostByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != 1 {
		t.Fatalf("expected post ID 1, got %d", post.ID)
	}
}

func TestGetPostByID_Error(t *testing.T) {
	repo, mock := setupTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE deleted_at IS NULL AND "posts"."id" = $1 ORDER BY "posts"."id" LIMIT $2`)).WithArgs(1, 1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetPostByID(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdatePost_Success(t *testing.T) {
	repo, mock := setupTestDB(t)
	now := time.Now()
	post := entity.Post{ID: 1, Title: "Updated", Content: "Content", UserID: 1, ImageURL: "", CreatedAt: now, UpdatedAt: now, DeletedAt: nil}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts" SET`)).WithArgs(
		post.ID, post.UserID, post.Title, post.Content, sqlmock.AnyArg(), sqlmock.AnyArg(), post.ID,
	).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.UpdatePost(ctx, post)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdatePost_Error(t *testing.T) {
	repo, mock := setupTestDB(t)
	now := time.Now()
	post := entity.Post{ID: 1, Title: "Updated", Content: "Content", UserID: 1, ImageURL: "", CreatedAt: now, UpdatedAt: now, DeletedAt: nil}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts" SET`)).WithArgs(
		post.ID, post.UserID, post.Title, post.Content, sqlmock.AnyArg(), sqlmock.AnyArg(), post.ID,
	).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.UpdatePost(ctx, post)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDeletePost_Success(t *testing.T) {
	repo, mock := setupTestDB(t)
	post := entity.Post{ID: 1, Title: "ToDelete", Content: "Content", UserID: 1}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts" SET`)).WithArgs(post.UserID, post.Title, post.Content, sqlmock.AnyArg(), post.ID, post.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	ctx := context.Background()
	err := repo.DeletePost(ctx, post)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeletePost_Error(t *testing.T) {
	repo, mock := setupTestDB(t)
	post := entity.Post{ID: 1, Title: "ToDelete", Content: "Content", UserID: 1}
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts" SET`)).WithArgs(post.UserID, post.Title, post.Content, sqlmock.AnyArg(), post.ID, post.ID).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	err := repo.DeletePost(ctx, post)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSearchPosts_Success(t *testing.T) {
	repo, mock := setupTestDB(t)
	userID := 1
	query := "test"
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "posts" WHERE deleted_at IS NULL AND (user_id = $1 OR is_public = true) AND (title ILIKE $2 OR content ILIKE $3)`)).
		WithArgs(userID, "%test%", "%test%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE deleted_at IS NULL AND (user_id = $1 OR is_public = true) AND (title ILIKE $2 OR content ILIKE $3) ORDER BY created_at DESC LIMIT $4`)).
		WithArgs(userID, "%test%", "%test%", 10).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "image_url", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, "Test", "Content", "url", now, now, nil))
	ctx := context.Background()
	posts, total, err := repo.SearchPosts(ctx, userID, query, 0, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
}

func TestSearchPosts_Error(t *testing.T) {
	repo, mock := setupTestDB(t)
	query := "test"
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "posts" WHERE deleted_at IS NULL AND (title ILIKE $1 OR content ILIKE $2)`)).WithArgs("%test%", "%test%").WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, _, err := repo.SearchPosts(ctx, 1, query, 0, 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
