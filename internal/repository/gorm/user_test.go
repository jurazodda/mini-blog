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

func setupUserTestDB(t *testing.T) (*GormRepository, sqlmock.Sqlmock) {
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

func TestCreateUser_Success(t *testing.T) {
	repo, mock := setupUserTestDB(t)
	user := entity.User{ID: 1, Username: "test", Email: "test@example.com"}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).WithArgs(
		user.Username, user.Email, user.PasswordHash, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.ID,
	).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	ctx := context.Background()
	_, err := repo.CreateUser(ctx, user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateUser_Error(t *testing.T) {
	repo, mock := setupUserTestDB(t)
	user := entity.User{ID: 1, Username: "test", Email: "test@example.com"}
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).WithArgs(
		user.Username, user.Email, user.PasswordHash, sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.ID,
	).WillReturnError(gorm.ErrInvalidDB)
	mock.ExpectRollback()
	ctx := context.Background()
	_, err := repo.CreateUser(ctx, user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUserByEmail_Success(t *testing.T) {
	repo, mock := setupUserTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND deleted_at IS NULL ORDER BY "users"."id" LIMIT $2`)).WithArgs("test@example.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at", "deleted_at"}).AddRow(1, "test", "test@example.com", "hash", now, now, nil))
	ctx := context.Background()
	user, err := repo.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 1 {
		t.Fatalf("expected user ID 1, got %d", user.ID)
	}
}

func TestGetUserByEmail_Error(t *testing.T) {
	repo, mock := setupUserTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).WithArgs("test@example.com", 1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetUserByEmail(ctx, "test@example.com")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetUserByID_Success(t *testing.T) {
	repo, mock := setupUserTestDB(t)
	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE deleted_at IS NULL AND "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at", "deleted_at"}).AddRow(1, "test", "test@example.com", "hash", now, now, nil))
	ctx := context.Background()
	user, err := repo.GetUserByID(ctx, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 1 {
		t.Fatalf("expected user ID 1, got %d", user.ID)
	}
}

func TestGetUserByID_Error(t *testing.T) {
	repo, mock := setupUserTestDB(t)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE deleted_at IS NULL AND "users"."id" = $1 ORDER BY "users"."id" LIMIT $2`)).WithArgs(1, 1).WillReturnError(gorm.ErrInvalidDB)
	ctx := context.Background()
	_, err := repo.GetUserByID(ctx, 1)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
