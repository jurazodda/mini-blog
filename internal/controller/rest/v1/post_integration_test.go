package v1

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	gormrepo "mini-blog/internal/repository/gorm"
	"mini-blog/internal/service"
	"mini-blog/pkg/logger"

	"mini-blog/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestServerWithSqlMock(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, int, string) {
	// Установим SigningKey для JWT
	viper.Set("jwt.signingkey", "test-signing-key")
	t.Setenv("JWT_SIGNINGKEY", "test-signing-key")

	// Ручная инициализация config
	config.Set(&config.Config{
		Jwt: config.Jwt{
			TokenTTL: 12,
		},
	})

	// Создаём sqlmock
	dbSQL, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	// Создаём *gorm.DB поверх sqlmock
	dialector := postgres.New(postgres.Config{Conn: dbSQL, PreferSimpleProtocol: true})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	log := logger.NewTestLogger()
	repo := gormrepo.NewGormRepository(db, log)
	svc := service.NewService(repo, log)

	h := &RestHandlerv1{Service: svc, Logger: log}
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h.InitRoutes(r)

	userID := 42
	ctx := context.Background()
	token, err := service.GenerateToken(ctx, userID)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	return r, mock, userID, token
}

func TestCreatePost_Integration(t *testing.T) {
	router, mock, userID, token := setupTestServerWithSqlMock(t)

	// Ожидание на sqlmock: вставка поста
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO \"posts\"").
		WithArgs(
			userID,           // user_id
			"Test Title",     // title
			"",               // image_url
			"Test Content",   // content
			false,            // is_public (по умолчанию)
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			nil,              // deleted_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Формируем HTTP-запрос
	body := map[string]interface{}{
		"title":   "Test Title",
		"content": "Test Content",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreatePost_ValidationError(t *testing.T) {
	router, _, _, token := setupTestServerWithSqlMock(t)
	body := map[string]interface{}{
		"title":   "",
		"content": "Test Content",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreatePost_DBError(t *testing.T) {
	router, mock, userID, token := setupTestServerWithSqlMock(t)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "posts"`).
		WithArgs(userID, "Test Title", "", "Test Content", sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()
	body := map[string]interface{}{
		"title":   "Test Title",
		"content": "Test Content",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/post", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetPosts_Success(t *testing.T) {
	router, mock, _, token := setupTestServerWithSqlMock(t)
	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL.*`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "image_url", "content", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 42, "Test Title", "", "Test Content", now, now, nil))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "comments".*`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "likes".*`).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	req := httptest.NewRequest(http.MethodGet, "/post", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetPosts_DBError(t *testing.T) {
	router, mock, _, token := setupTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL.*`).WillReturnError(sql.ErrConnDone)
	req := httptest.NewRequest(http.MethodGet, "/post", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetPostByID_Success(t *testing.T) {
	router, mock, userID, token := setupTestServerWithSqlMock(t)
	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1.*`).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "image_url", "content", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, userID, "Test Title", "", "Test Content", now, now, nil))
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE post_id = \$1 AND deleted_at IS NULL`).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "created_at", "updated_at", "deleted_at"}))
	mock.ExpectQuery(`SELECT \* FROM "likes" WHERE post_id = \$1 AND deleted_at IS NULL`).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id", "created_at", "updated_at", "deleted_at"}))
	req := httptest.NewRequest(http.MethodGet, "/post/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetPostByID_NotFound(t *testing.T) {
	router, mock, _, token := setupTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1.*`).WithArgs(1, 1).WillReturnError(sql.ErrNoRows)
	req := httptest.NewRequest(http.MethodGet, "/post/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 404 or 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetPostByID_DBError(t *testing.T) {
	router, mock, _, token := setupTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1.*`).WithArgs(1, 1).WillReturnError(sql.ErrConnDone)
	req := httptest.NewRequest(http.MethodGet, "/post/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdatePost_Success(t *testing.T) {
	router, mock, userID, token := setupTestServerWithSqlMock(t)
	now := time.Now()
	// Ожидание: select для проверки владельца
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1.*`).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "image_url", "content", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, userID, "Old Title", "", "Old Content", now, now, nil))
	// update
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "posts" SET`).
		WithArgs(1, userID, "New Title", "New Content", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := map[string]interface{}{
		"title":   "New Title",
		"content": "New Content",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/post/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestDeletePost_Success(t *testing.T) {
	router, mock, userID, token := setupTestServerWithSqlMock(t)
	now := time.Now()
	// Ожидание: select для проверки владельца
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1.*`).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "image_url", "content", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, userID, "Title", "", "Content", now, now, nil))
	// Ожидание: update (только реально сериализуемые GORM поля для soft delete)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "posts" SET`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodDelete, "/post/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestSearchPosts_Success(t *testing.T) {
	router, mock, userID, token := setupTestServerWithSqlMock(t)
	now := time.Now()
	mock.ExpectQuery(`SELECT count\(\*\) FROM "posts" WHERE deleted_at IS NULL AND \(user_id = \$1 OR is_public = true\) AND \(title ILIKE \$2 OR content ILIKE \$3\)`).
		WithArgs(userID, "%test%", "%test%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND \(user_id = \$1 OR is_public = true\) AND \(title ILIKE \$2 OR content ILIKE \$3\) ORDER BY created_at DESC LIMIT \$4`).
		WithArgs(userID, "%test%", "%test%", 10).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "image_url", "content", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, userID, "Test Title", "", "Test Content", now, now, nil))

	req := httptest.NewRequest(http.MethodGet, "/post/search?q=test&page=1&page_size=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Logf("body: %s", w.Body.String())
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
