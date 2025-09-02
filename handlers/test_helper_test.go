package handlers

import (
	"io"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const secretString = "secret"

func init() {
	middlewares.SetSecret(func() []byte { return []byte(secretString) })
}

func setupMockGorm(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, func()) {
	t.Helper()
	sqlDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	gdb, _ := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		},
	)
	cleanup := func() { _ = sqlDB.Close() }
	return mock, gdb, cleanup
}

func setupApp(gdb *gorm.DB) *fiber.App {
	app := fiber.New()
	// inject mocked DB into request context
	app.Use(middlewares.DBMiddleware(gdb))
	// now mount routes
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

func makeJWT(t *testing.T, secret []byte, userID uint) string {
	t.Helper()

	claims := jwt.MapClaims{
		"user_id": userID,
		"iat":     time.Now().UTC().Unix(),
		"exp":     time.Now().UTC().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return signed
}

func preloadUserForAuth(mock sqlmock.Sqlmock, userID uint, hasAdmin bool, hasTeacher bool, hasLearner bool) {
	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"\."id" = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT .*`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))
	if hasAdmin {
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(99, userID))
	} else {
		mock.ExpectQuery(`SELECT \* FROM "admins" WHERE "admins"\."user_id" = \$1 AND "admins"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}
	if hasTeacher {
		mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."user_id" = \$1 AND "teachers"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(99, userID))
	} else {
		mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."user_id" = \$1 AND "teachers"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}
	if hasLearner {
		mock.ExpectQuery(`SELECT \* FROM "learners" WHERE "learners"\."user_id" = \$1 AND "learners"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(99, userID))
	} else {
		mock.ExpectQuery(`SELECT \* FROM "learners" WHERE "learners"\."user_id" = \$1 AND "learners"\."deleted_at" IS NULL`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}

}
