package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

func setupRepostTestServerWithSqlMock(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, int, string) {
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

// Пример: создание репоста
func TestCreateRepost_Success(t *testing.T) {
	router, mock, userID, token := setupRepostTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1.*`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "image_url", "created_at", "updated_at", "deleted_at"}).AddRow(1, userID, "title", "content", "url", nil, nil, nil))
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "reposts"`).WithArgs(userID, 1, "Repost title", "Repost content", "url", 0, 0, sqlmock.AnyArg(), nil).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	body := map[string]interface{}{"post_id": 1, "title": "Repost title", "content": "Repost content"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/repost", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetReposts_Success(t *testing.T) {
	router, mock, _, token := setupRepostTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "reposts" WHERE deleted_at IS NULL.*`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at", "image_url"}).
			AddRow(1, 1, 42, "title", "content", nil, nil, nil, "url"))
	req := httptest.NewRequest(http.MethodGet, "/repost", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetRepostByID_Success(t *testing.T) {
	router, mock, userID, token := setupRepostTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "reposts" WHERE deleted_at IS NULL AND "reposts"\."id" = \$1.*`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at", "image_url"}).AddRow(1, 1, userID, "title", "content", nil, nil, nil, "url"))
	req := httptest.NewRequest(http.MethodGet, "/repost/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestUpdateRepost_Success(t *testing.T) {
	router, mock, userID, token := setupRepostTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "reposts" WHERE deleted_at IS NULL AND "reposts"\."id" = \$1 ORDER BY "reposts"\."id" LIMIT \$2`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at", "image_url"}).AddRow(1, 1, userID, "old title", "old content", nil, nil, nil, "url"))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "reposts" SET`).WithArgs(1, userID, "Updated title", "Updated content", 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	body := map[string]interface{}{"title": "Updated title", "content": "Updated content"}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPut, "/repost/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestDeleteRepost_Success(t *testing.T) {
	router, mock, userID, token := setupRepostTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "reposts" WHERE deleted_at IS NULL AND "reposts"\."id" = \$1.*`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at", "image_url"}).AddRow(1, 1, userID, "title", "content", nil, nil, nil, "url"))
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "reposts" SET`).WithArgs(sqlmock.AnyArg(), 1, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	req := httptest.NewRequest(http.MethodDelete, "/repost/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetRepostsByUserID_Success(t *testing.T) {
	router, mock, _, token := setupRepostTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "reposts" WHERE user_id = \$1 AND deleted_at IS NULL`).WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "title", "content", "created_at", "updated_at", "deleted_at", "image_url"}).AddRow(1, 1, 42, "title", "content", nil, nil, nil, "url"))
	req := httptest.NewRequest(http.MethodGet, "/repost/user/42", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

// ... Аналогично реализовать тесты на ошибки и edge-cases ...
