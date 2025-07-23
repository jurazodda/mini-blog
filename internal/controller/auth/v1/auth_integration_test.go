package v1

import (
	"bytes"
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

func setupAuthTestServerWithSqlMock(t *testing.T) *gin.Engine {
	// Подготавливаем тестовый сервер с sqlmock
	viper.Set("jwt.signingkey", "test-signing-key")
	t.Setenv("JWT_SIGNINGKEY", "test-signing-key")
	config.Set(&config.Config{
		Jwt: config.Jwt{
			TokenTTL: 12,
		},
	})

	dbSQL, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{Conn: dbSQL, PreferSimpleProtocol: true})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm DB: %v", err)
	}

	log, _ := logger.NewLogger()
	repo := gormrepo.NewGormRepository(db, log)
	svc := service.NewService(repo, log)
	h := &AuthHandlerV1{
		Service: svc,
		Logger:  log,
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	h.InitRoutes(r)
	
	return r
}

func TestSignUp_Success(t *testing.T) {
	// Test-case: успешное регистрация пользователя
	r, mock := setupAuthTestServerWithSqlMockAndReturnMock(t)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WithArgs("testuser", "test@example.com", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	
	body := map[string]interface{}{
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	
	req := httptest.NewRequest(http.MethodPost, "/auth/v1/sign-up", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestSignUp_ValidationError(t *testing.T) {
	// Test-case: ошибка валидации при регистрации пользователя
	r := setupAuthTestServerWithSqlMock(t)
	
	body := map[string]interface{}{
		"username": "",
		"email":    "bad",
		"password": "",
	}
	jsonBody, _ := json.Marshal(body)
	
	req := httptest.NewRequest(http.MethodPost, "/auth/v1/sign-up", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 400 or 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func setupAuthTestServerWithSqlMockAndReturnMock(t *testing.T) (*gin.Engine, sqlmock.Sqlmock) {
	// Подготавливаем тестовый сервер с sqlmock
	viper.Set("jwt.signingkey", "test-signing-key")
	t.Setenv("JWT_SIGNINGKEY", "test-signing-key")
	config.Set(&config.Config{
		Jwt: config.Jwt{
			TokenTTL: 12,
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

	log, _ := logger.NewLogger()
	repo := gormrepo.NewGormRepository(db, log)
	svc := service.NewService(repo, log)
	h := &AuthHandlerV1{
		Service: svc,
		Logger:  log,
	}
	
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h.InitRoutes(r)
	
	return r, mock
}

func TestSignIn_Success(t *testing.T) {
	// Test-case: успешное вход в систему
	r, mock := setupAuthTestServerWithSqlMockAndReturnMock(t)
	
	now := time.Now()
	const hash = "$2a$10$HEprfLqaBWP9rUsuQtq4.u0YZo7QanIXuOjLA3Wx95PwGeYZhobH2"
	
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND deleted_at IS NULL ORDER BY "users"\."id" LIMIT \$2`).WithArgs("test@example.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at", "deleted_at"}).
			AddRow(1, "testuser", "test@example.com", hash, now, now, nil))

	body := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	
	req := httptest.NewRequest(http.MethodPost, "/auth/v1/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}

	resp := make(map[string]interface{})

	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == nil {
		t.Fatalf("expected token in response, got: %s", w.Body.String())
	}
	
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestSignIn_ValidationError(t *testing.T) {
	// Test-case: ошибка валидации при входе в систему
	r := setupAuthTestServerWithSqlMock(t)
	
	body := map[string]interface{}{
		"email":    "",
		"password": "",
	}
	jsonBody, _ := json.Marshal(body)
	
	req := httptest.NewRequest(http.MethodPost, "/auth/v1/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 400 or 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestSignIn_UserNotFound(t *testing.T) {
	// Test-case: пользователь не найден
	r, mock := setupAuthTestServerWithSqlMockAndReturnMock(t)
	
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("notfound@example.com", 1).
		WillReturnError(sql.ErrNoRows)
	
	body := map[string]interface{}{
		"email":    "notfound@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	
	req := httptest.NewRequest(http.MethodPost, "/auth/v1/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestSignIn_InvalidPassword(t *testing.T) {
	// Test-case: неверный пароль
	r, mock := setupAuthTestServerWithSqlMockAndReturnMock(t)
	
	now := time.Now()
	
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs("test@example.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at", "deleted_at"}).
		AddRow(1, "testuser", "test@example.com", "$2a$10$invalidhash", now, now, nil))
	
	body := map[string]interface{}{
		"email":    "test@example.com",
		"password": "wrongpassword",
	}
	jsonBody, _ := json.Marshal(body)
	
	req := httptest.NewRequest(http.MethodPost, "/auth/v1/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestSignIn_DBError(t *testing.T) {
	// Test-case: ошибка при запросе к базе данных
	r, mock := setupAuthTestServerWithSqlMockAndReturnMock(t)
	
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 ORDER BY "users"\."id" LIMIT \$2`).
	WithArgs("test@example.com", 1).
	WillReturnError(sql.ErrConnDone)
	
	body := map[string]interface{}{
		"email":    "test@example.com",
		"password": "password123",
	}
	jsonBody, _ := json.Marshal(body)
	
	req := httptest.NewRequest(http.MethodPost, "/auth/v1/sign-in", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d, body: %s", w.Code, w.Body.String())
	}
}
