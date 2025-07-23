package v1

import (
	"bytes"
	"context"
	"database/sql"
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

func setupLikeTestServerWithSqlMock(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, int, string) {
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

func TestCreateLikePost_Success(t *testing.T) {
	router, mock, userID, token := setupLikeTestServerWithSqlMock(t)
	// Проверка, что пост существует
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1 ORDER BY "posts"\."id" LIMIT \$2`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "image_url", "created_at", "updated_at", "deleted_at"}).AddRow(1, userID, "title", "content", "url", nil, nil, nil))
	// Нет лайка, создаём
	mock.ExpectQuery(`SELECT \* FROM "likes" WHERE user_id = \$1 AND post_id = \$2 AND deleted_at IS NULL ORDER BY "likes"."id" LIMIT \$3`).WithArgs(userID, 1, sqlmock.AnyArg()).WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "likes"`).WithArgs(userID, 1, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	body := map[string]interface{}{"post_id": 1}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/like/post", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Log(w.Body.String())
		t.Logf("mock expectations: %+v", mock)
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Log(err)
		t.Fatalf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreateLikePost_RemoveLike(t *testing.T) {
	router, mock, userID, token := setupLikeTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1.*`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "title", "content", "image_url", "created_at", "updated_at", "deleted_at"}).AddRow(1, userID, "title", "content", "url", nil, nil, nil))
	mock.ExpectQuery(`SELECT \* FROM "likes" WHERE user_id = \$1 AND post_id = \$2 AND deleted_at IS NULL ORDER BY "likes"."id" LIMIT \$3`).WithArgs(userID, 1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "post_id"}).AddRow(1, userID, 1))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "likes" WHERE`).WithArgs(userID, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	body := map[string]interface{}{"post_id": 1}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/like/post", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateLikePost_DBError(t *testing.T) {
	router, mock, _, token := setupLikeTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "posts" WHERE deleted_at IS NULL AND "posts"\."id" = \$1.*`).WithArgs(1, 1).WillReturnError(sql.ErrConnDone)
	body := map[string]interface{}{"post_id": 1}
	jsonBody, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/like/post", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateLikeComment_Success(t *testing.T) {
	router, mock, _, token := setupLikeTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"."id" LIMIT \$2`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, 42, "comment", nil, nil, nil))
	mock.ExpectQuery(`SELECT \* FROM "comment_likes" WHERE user_id = \$1 AND comment_id = \$2 ORDER BY "comment_likes"\."id" LIMIT \$3`).WithArgs(42, 1, sqlmock.AnyArg()).WillReturnError(gorm.ErrRecordNotFound)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "comment_likes"`).WithArgs(42, 1, sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	req := httptest.NewRequest(http.MethodPost, "/like/comment/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateLikeComment_RemoveLike(t *testing.T) {
	router, mock, _, token := setupLikeTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"."id" LIMIT \$2`).WithArgs(1, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "post_id", "user_id", "comment", "created_at", "updated_at", "deleted_at"}).AddRow(1, 1, 42, "comment", nil, nil, nil))
	mock.ExpectQuery(`SELECT \* FROM "comment_likes" WHERE user_id = \$1 AND comment_id = \$2 ORDER BY "comment_likes"\."id" LIMIT \$3`).WithArgs(42, 1, sqlmock.AnyArg()).WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "comment_id"}).AddRow(1, 42, 1))
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "comment_likes" WHERE`).WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	req := httptest.NewRequest(http.MethodPost, "/like/comment/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestCreateLikeComment_DBError(t *testing.T) {
	router, mock, _, token := setupLikeTestServerWithSqlMock(t)
	mock.ExpectQuery(`SELECT \* FROM "comments" WHERE deleted_at IS NULL AND "comments"\."id" = \$1 ORDER BY "comments"."id" LIMIT \$2`).WithArgs(1, 1).WillReturnError(sql.ErrConnDone)
	req := httptest.NewRequest(http.MethodPost, "/like/comment/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}
