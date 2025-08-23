package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockGorm(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, func()) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	cleanup := func() { _ = sqlDB.Close() }
	return mock, gdb, cleanup
}

func init() {
	SetSecretProvider(func() []byte { return []byte("secret") })
}

func makeJWT(t *testing.T, userID uint) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte("secret"))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return s
}

func preloadUserForAuth(mock sqlmock.Sqlmock, userID uint, hasAdmin, hasTeacher, hasLearner bool) {
	mock.MatchExpectationsInOrder(false)

	// users.First (LIMIT param is bound by GORM -> WithArgs(userID, 1))
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

	// Admin preload
	if hasAdmin {
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(10, userID))
	} else {
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}

	// Learner preload
	if hasLearner {
		mock.ExpectQuery(`SELECT \* FROM "learners" WHERE "learners"\."user_id" = \$1 AND "learners"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(20, userID))
	} else {
		mock.ExpectQuery(`SELECT \* FROM "learners" WHERE "learners"\."user_id" = \$1 AND "learners"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}

	// Teacher preload
	if hasTeacher {
		mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."user_id" = \$1 AND "teachers"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(30, userID))
	} else {
		mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."user_id" = \$1 AND "teachers"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}
}

/* --------------------- tests --------------------- */

func TestProtectedMiddleware_MissingToken_401(t *testing.T) {
	_, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/secure", ProtectedMiddleware(), func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	// no Authorization header
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestProtectedMiddleware_InvalidToken_401(t *testing.T) {
	_, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/secure", ProtectedMiddleware(), func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestProtectedMiddleware_Valid_200(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(42)
	preloadUserForAuth(mock, userID, false, false, false)

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/secure", ProtectedMiddleware(), func(c *fiber.Ctx) error { return c.SendStatus(200) })

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, userID))
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAdminRequired_Success_200(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(7)
	preloadUserForAuth(mock, userID, true, false, false) // has Admin

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/admin-only", ProtectedMiddleware(), AdminRequired(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, userID))
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestAdminRequired_Forbidden_403(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(8)
	preloadUserForAuth(mock, userID, false, false, false) // no Admin

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/admin-only", ProtectedMiddleware(), AdminRequired(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, userID))
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusForbidden)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestTeacherRequired_Success_200(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(9)
	preloadUserForAuth(mock, userID, false, true, false) // has Teacher

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/teacher-only", ProtectedMiddleware(), TeacherRequired(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/teacher-only", nil)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, userID))
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestTeacherRequired_Forbidden_403(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(9)
	preloadUserForAuth(mock, userID, false, false, false) // no Teacher

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/teacher-only", ProtectedMiddleware(), TeacherRequired(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/teacher-only", nil)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, userID))
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusForbidden)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestLearnerRequired_Success_200(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(10)
	preloadUserForAuth(mock, userID, false, false, true) // has Learner

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/learner-only", ProtectedMiddleware(), LearnerRequired(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/learner-only", nil)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, userID))
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusOK)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestLearnerRequired_Forbidden_403(t *testing.T) {
	mock, gdb, cleanup := setupMockGorm(t)
	defer cleanup()

	userID := uint(9)
	preloadUserForAuth(mock, userID, false, false, false) // no Learner

	app := fiber.New()
	app.Use(DBMiddleware(gdb))
	app.Get("/learner-only", ProtectedMiddleware(), LearnerRequired(), func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	req := httptest.NewRequest(http.MethodGet, "/learner-only", nil)
	req.Header.Set("Authorization", "Bearer "+makeJWT(t, userID))
	resp, _ := app.Test(req, -1)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status=%d want=%d", resp.StatusCode, http.StatusForbidden)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}
