package handlers

import (
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/config/dbserver"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockGorm(t *testing.T) (sqlmock.Sqlmock, func()) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	gdb, err := gorm.Open(postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	// Point global DB used by middleware/handlers
	dbserver.DB = gdb

	cleanup := func() { _ = sqlDB.Close() }
	return mock, cleanup
}

func setupApp() *fiber.App {
	app := fiber.New()
	AllRoutes(app)
	return app
}

func readBody(t *testing.T, r io.Reader) []byte {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return b
}

func makeJWT(t *testing.T, userID uint) string {
	t.Helper()
	if err := godotenv.Load("../.env"); err != nil {
		log.Println(".env file not found, using system environment variables")
	}
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	secret := []byte(secretStr)
	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(1 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return s
}

func preloadUserForAuth(mock sqlmock.Sqlmock, userID uint) {
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
		WillReturnRows(sqlmock.NewRows([]string{"id"}))
}
