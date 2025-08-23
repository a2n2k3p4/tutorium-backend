package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/config/dbserver"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockGormForMW(t *testing.T) sqlmock.Sqlmock {
	t.Helper()
	jwtSecret = []byte("secret")

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })

	gdb, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}),
		&gorm.Config{},
	)
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}

	dbserver.DB = gdb
	return mock
}

func makeJWTForMW(t *testing.T, secret []byte, userID uint) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString(secret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

// 401 when Authorization header missing
func TestProtected_NoToken_401(t *testing.T) {
	_ = setupMockGormForMW(t)

	app := fiber.New()
	app.Get("/x", ProtectedMiddleware(), func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

// 401 when token is wrong
func TestProtected_BadToken_401(t *testing.T) {
	_ = setupMockGormForMW(t)

	app := fiber.New()
	app.Get("/x", ProtectedMiddleware(), func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

// 401 when token valid but user not found
func TestProtected_UserNotFound_401(t *testing.T) {
	mock := setupMockGormForMW(t)

	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(42, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	app := fiber.New()
	app.Get("/x", ProtectedMiddleware(), func(c *fiber.Ctx) error { return c.SendStatus(200) })

	token := makeJWTForMW(t, []byte("secret"), 42)
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 403 when user exists but has no Admin role
func TestAdminRequired_Forbidden_403(t *testing.T) {
	mock := setupMockGormForMW(t)

	userID := 42
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "learners" .*`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "teachers" .*`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	app := fiber.New()
	app.Get("/only-admin",
		ProtectedMiddleware(),
		AdminRequired(),
		func(c *fiber.Ctx) error { return c.SendStatus(200) },
	)

	token := makeJWTForMW(t, []byte("secret"), uint(userID))
	req := httptest.NewRequest(http.MethodGet, "/only-admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

// 200 when user exists and has Admin role
func TestAdminRequired_Success_200(t *testing.T) {
	mock := setupMockGormForMW(t)

	userID := 42
	adminID := 90
	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(adminID, userID))
	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE "learners"\."user_id" = \$1 AND "learners"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."user_id" = \$1 AND "teachers"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	app := fiber.New()
	app.Get("/only-admin",
		ProtectedMiddleware(),
		AdminRequired(),
		func(c *fiber.Ctx) error { return c.SendStatus(200) },
	)

	token := makeJWTForMW(t, []byte("secret"), uint(userID))
	req := httptest.NewRequest(http.MethodGet, "/only-admin", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestTeacherRequired_Forbidden_403(t *testing.T) {
	mock := setupMockGormForMW(t)

	userID := 42
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "learners" .*`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "teachers" .*`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	app := fiber.New()
	app.Get("/only-teacher",
		ProtectedMiddleware(),
		TeacherRequired(),
		func(c *fiber.Ctx) error { return c.SendStatus(200) },
	)

	token := makeJWTForMW(t, []byte("secret"), uint(userID))
	req := httptest.NewRequest(http.MethodGet, "/only-teacher", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
func TestTeacherRequired_Success_200(t *testing.T) {
	mock := setupMockGormForMW(t)

	userID := 42
	teacherID := 90
	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE "learners"\."user_id" = \$1 AND "learners"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."user_id" = \$1 AND "teachers"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(teacherID, userID))

	app := fiber.New()
	app.Get("/only-teacher",
		ProtectedMiddleware(),
		TeacherRequired(),
		func(c *fiber.Ctx) error { return c.SendStatus(200) },
	)

	token := makeJWTForMW(t, []byte("secret"), uint(userID))
	req := httptest.NewRequest(http.MethodGet, "/only-teacher", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestLearnerRequired_Forbidden_403(t *testing.T) {
	mock := setupMockGormForMW(t)

	userID := 42
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "learners" .*`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "teachers" .*`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	app := fiber.New()
	app.Get("/only-learner",
		ProtectedMiddleware(),
		LearnerRequired(),
		func(c *fiber.Ctx) error { return c.SendStatus(200) },
	)

	token := makeJWTForMW(t, []byte("secret"), uint(userID))
	req := httptest.NewRequest(http.MethodGet, "/only-learner", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusForbidden)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
func TestLearnerRequired_Success_200(t *testing.T) {
	mock := setupMockGormForMW(t)

	userID := 42
	learnerID := 90
	mock.MatchExpectationsInOrder(false)
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	mock.ExpectQuery(`SELECT \* FROM "learners" WHERE "learners"\."user_id" = \$1 AND "learners"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(learnerID, userID))
	mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."user_id" = \$1 AND "teachers"\."deleted_at" IS NULL`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
	app := fiber.New()
	app.Get("/only-learner",
		ProtectedMiddleware(),
		LearnerRequired(),
		func(c *fiber.Ctx) error { return c.SendStatus(200) },
	)

	token := makeJWTForMW(t, []byte("secret"), uint(userID))
	req := httptest.NewRequest(http.MethodGet, "/only-learner", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
