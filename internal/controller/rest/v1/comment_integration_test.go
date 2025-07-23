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

	"mini-blog/config"
	gormrepo "mini-blog/internal/repository/gorm"
	"mini-blog/internal/service"
	"mini-blog/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupCommentTestServerWithSqlMock(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, int, string) {
	viper.Set("jwt.signingkey", "test-signing-key")
	t.Setenv("JWT_SIGNINGKEY", "test-signing-key")
	config.Set(&config.Config{
		Jwt: config.Jwt{
			TokenTTL:   12,
		},
	})
	dbSQL, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
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

func TestCreateComment_Integration(t *testing.T) {
	router, mock, userID, token := setupCommentTestServerWithSqlMock(t)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "comments"`).
		WithArgs(userID, 1, "Test comment", 0, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	body := map[string]interface{}{
		"post_id": 1,
		"comment": "Test comment",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/comment", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateComment_ValidationError(t *testing.T) {
	router, _, _, token := setupCommentTestServerWithSqlMock(t)
	body := map[string]interface{}{
		"post_id": 0,
		"comment": "",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/comment", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateComment_DBError(t *testing.T) {
	router, mock, userID, token := setupCommentTestServerWithSqlMock(t)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "comments"`).
		WithArgs(userID, 1, "Test comment", 0, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnError(sql.ErrConnDone)
	mock.ExpectRollback()
	body := map[string]interface{}{
		"post_id": 1,
		"comment": "Test comment",
	}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/comment", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCommentByID_Success(t *testing.T) {
	router, mock, userID, token := setupCommentTestServerWithSqlMock(t)
	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"\."id" LIMIT \$2`).WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 1, userID, "Test comment", now, now, nil))
	req := httptest.NewRequest(http.MethodGet, "/comment/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCommentByID_NotFound(t *testing.T) {
	router, mock, _, token := setupCommentTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"\."id" LIMIT \$2`).WithArgs(1, 1).WillReturnError(sql.ErrNoRows)
	req := httptest.NewRequest(http.MethodGet, "/comment/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 404 or 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCommentByID_DBError(t *testing.T) {
	router, mock, _, token := setupCommentTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"\."id" LIMIT \$2`).WithArgs(1, 1).WillReturnError(sql.ErrConnDone)
	req := httptest.NewRequest(http.MethodGet, "/comment/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateComment_Success(t *testing.T) {
	router, mock, userID, token := setupCommentTestServerWithSqlMock(t)
	now := time.Now()
	// Ожидаем получение существующего комментария
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"\."id" LIMIT \$2`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, userID, "Old comment", now, now, nil))
	// Ожидаем UPDATE
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "comments" SET`).WithArgs(1, userID, "Updated comment", sqlmock.AnyArg(), 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	body := map[string]interface{}{"comment": "Updated comment"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/comment/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateComment_ValidationError(t *testing.T) {
	router, _, _, token := setupCommentTestServerWithSqlMock(t)
	body := map[string]interface{}{"comment": ""}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/comment/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateComment_NotFound(t *testing.T) {
	router, mock, _, token := setupCommentTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"\."id" LIMIT \$2`).WithArgs(1, 1).WillReturnError(sql.ErrNoRows)
	body := map[string]interface{}{"comment": "Updated comment"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/comment/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 404 or 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteComment_Success(t *testing.T) {
	router, mock, userID, token := setupCommentTestServerWithSqlMock(t)
	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"\."id" LIMIT \$2`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, userID, "Test comment", now, now, nil))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "comments" SET`).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	req := httptest.NewRequest(http.MethodDelete, "/comment/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteComment_NotFound(t *testing.T) {
	router, mock, _, token := setupCommentTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"\."id" LIMIT \$2`).WithArgs(1, 1).WillReturnError(sql.ErrNoRows)
	req := httptest.NewRequest(http.MethodDelete, "/comment/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 404 or 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCommentsByPostID_Success(t *testing.T) {
	router, mock, _, token := setupCommentTestServerWithSqlMock(t)
	now := time.Now()
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE post_id = \$1 AND deleted_at IS NULL`).WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, 1, 42, "Test comment", now, now, nil))
	req := httptest.NewRequest(http.MethodGet, "/comment/post/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetCommentsByPostID_DBError(t *testing.T) {
	router, mock, _, token := setupCommentTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE post_id = \$1 AND deleted_at IS NULL`).WithArgs(1).WillReturnError(sql.ErrConnDone)
	req := httptest.NewRequest(http.MethodGet, "/comment/post/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}
