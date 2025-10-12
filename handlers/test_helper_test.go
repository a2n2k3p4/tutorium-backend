package handlers

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/config"
	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/valyala/fasthttp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

/* ------------------ SetSecret Helper ------------------ */
const secretString = "secret"

func init() {
	middlewares.SetSecret(func() []byte { return []byte(secretString) })
}

/* ------------------ Test set-up Helper  ------------------ */

func setupMockGorm(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, func()) {
	t.Helper()
	sqlDB, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	gdb, _ := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB, PreferSimpleProtocol: true}),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
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

/* ------------------ Reader Helper ------------------ */
func readBody(t *testing.T, r io.Reader) []byte {
	t.Helper()
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return b
}

func jsonBody(v any) []byte {
	b, _ := json.Marshal(v)
	return b
}

/* ------------------ JWT maker Helper ------------------ */
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

/* ------------------ Bypass test Helper ------------------ */
func RunInDifferentStatus(t *testing.T,
	body func(
		t *testing.T,
		mock sqlmock.Sqlmock,
		gdb *gorm.DB,
		app *fiber.App,
		payload *[]byte,
		uID *uint,
	),
	want int,
	method, path string,
) {
	t.Helper()
	cases := []struct {
		name string
		env  string
	}{
		{"bypass", "development"},
		{"unbypass", "production"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("STATUS", c.env)

			mock, gdb, cleanup := setupMockGorm(t)
			defer cleanup()
			mock.MatchExpectationsInOrder(false)

			app := setupApp(gdb)

			var payload []byte
			var uID uint = 0

			body(t, mock, gdb, app, &payload, &uID)

			input := httpInput{
				Method: method,
				Path:   path,
				UserID: &uID,
			}
			if len(payload) > 0 {
				input.Body = payload
				input.ContentType = "application/json"
			}
			resp := runHTTP(t, app, input)

			wantStatus(t, resp, want)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatalf("unmet expectations: %v", err)
			}
		})
	}
}

/* ------------------ Authentication Helper ------------------ */
func preloadUserForAuth(mock sqlmock.Sqlmock, userID uint, hasAdmin bool, hasTeacher bool, hasLearner bool) {
	if config.STATUS() == "development" {
		return
	}

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

/* ------------------ Expect query Helper ------------------ */
type Exp func(sqlmock.Sqlmock)

func ExpAuthUser(userID uint, asAdmin, asTeacher, asSomethingElse bool) Exp {
	return func(m sqlmock.Sqlmock) {
		preloadUserForAuth(m, userID, asAdmin, asTeacher, asSomethingElse)
	}
}

func ExpInsertReturningID(table string, id uint64) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin()
		m.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s".*RETURNING "id"`, table)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))
		m.ExpectCommit()
	}
}

func ExpDoubleInsertReturningID(table1 string, table2 string, id1 uint64, id2 uint64) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin()
		m.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s".*RETURNING "id"`, table1)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id1))
		m.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s".*RETURNING "id"`, table2)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id2))
		m.ExpectCommit()
	}
}

func ExpInsertError(table string, err error) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin()
		m.ExpectQuery(fmt.Sprintf(`INSERT INTO "%s".*RETURNING "id"`, table)).
			WillReturnError(err)
		m.ExpectRollback()
	}
}

func ExpListRows(table string, columns []string, rows ...[]any) Exp {
	return func(m sqlmock.Sqlmock) {
		r := sqlmock.NewRows(columns)
		for _, row := range rows {
			values := make([]driver.Value, len(row))
			for i, v := range row {
				values[i] = v
			}
			r.AddRow(values...)
		}
		m.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "%s".*`, table)).WillReturnRows(r)
	}
}

func ExpListError(table string, err error) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "%s".*`, table)).WillReturnError(err)
	}
}

func ExpSelectByIDFound(table string, id uint, cols []string, vals []any) Exp {
	return func(m sqlmock.Sqlmock) {
		values := make([]driver.Value, len(vals))
		for i, v := range vals {
			values[i] = v
		}
		m.ExpectQuery(fmt.Sprintf(
			`SELECT \* FROM "%s" WHERE id = \$1 AND "%s"\."deleted_at" IS NULL ORDER BY "%s"\."id" LIMIT .*`,
			table, table, table,
		)).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows(cols).AddRow(values...))
	}
}

func ExpSelectByIDEmpty(table string, id uint) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(fmt.Sprintf(
			`SELECT \* FROM "%s" WHERE id = \$1 AND "%s"\."deleted_at" IS NULL ORDER BY "%s"\."id" LIMIT .*`,
			table, table, table,
		)).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}
}

func ExpSelectByIDError(table string, id uint, err error) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(fmt.Sprintf(
			`SELECT \* FROM "%s" WHERE id = \$1 AND "%s"\."deleted_at" IS NULL ORDER BY "%s"\."id" LIMIT .*`,
			table, table, table,
		)).
			WithArgs(id, 1).
			WillReturnError(err)
	}
}

func ExpUpdateOK(table string) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin()
		m.ExpectExec(fmt.Sprintf(`UPDATE "%s" SET .* WHERE "%s"\."deleted_at" IS NULL`, table, table)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		m.ExpectCommit()
	}
}

func ExpUpdateError(table string, err error) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin()
		m.ExpectExec(fmt.Sprintf(`UPDATE "%s" SET .* WHERE "%s"\."deleted_at" IS NULL`, table, table)).
			WillReturnError(err)
		m.ExpectRollback()
	}
}

func ExpRecalculateClassRating(classID uint, avg float64) Exp {
	return func(m sqlmock.Sqlmock) {
		rows := sqlmock.NewRows([]string{"avg"}).AddRow(avg)
		m.ExpectQuery(`SELECT AVG\(rating\) FROM "reviews" WHERE class_id = \$1 AND "reviews"\."deleted_at" IS NULL`).
			WithArgs(classID).
			WillReturnRows(rows)

		m.ExpectBegin()
		m.ExpectExec(`UPDATE "classes" SET .* WHERE .*"classes"\."deleted_at" IS NULL`).
			WithArgs(avg, sqlmock.AnyArg(), classID).
			WillReturnResult(sqlmock.NewResult(0, 1))
		m.ExpectCommit()
	}
}

func ExpSoftDeleteOK(table string) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin()
		m.ExpectExec(fmt.Sprintf(`UPDATE "%s" SET "deleted_at"=`, table)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		m.ExpectCommit()
	}
}

func ExpSoftDeleteOKWithAllowNoTransaction(table string) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectExec(fmt.Sprintf(`UPDATE "%s" SET "deleted_at"=`, table)).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func ExpSoftDeleteError(table string, err error) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectBegin()
		m.ExpectExec(fmt.Sprintf(`UPDATE "%s" SET "deleted_at"=`, table)).
			WillReturnError(err)
		m.ExpectRollback()
	}
}

func ExpPreloadField(table string, columns []string, vals []any) Exp {
	return func(m sqlmock.Sqlmock) {
		values := make([]driver.Value, len(vals))
		for i, v := range vals {
			values[i] = v
		}
		r := sqlmock.NewRows(columns).AddRow(values...)
		m.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "%s".*`, table)).
			WillReturnRows(r)
	}
}

func ExpPreloadCanEmpty(table string, columns []string) Exp {
	return func(m sqlmock.Sqlmock) {
		m.ExpectQuery(fmt.Sprintf(`SELECT .* FROM "%s".*`, table)).
			WillReturnRows(sqlmock.NewRows(columns))
	}
}

func ExpPreloadM2M(joinTable string, childTable string, parentKey string,
	childKey string, pairs [][2]any) Exp {
	return func(m sqlmock.Sqlmock) {

		joinRows := sqlmock.NewRows([]string{parentKey, childKey})
		childIDs := make(map[any]struct{}, len(pairs))
		for _, p := range pairs {
			joinRows.AddRow(p[0], p[1])
			childIDs[p[1]] = struct{}{}
		}
		m.ExpectQuery(`SELECT .* FROM "` + joinTable + `".*WHERE .*"` + joinTable + `"\."` + parentKey + `" (=\s*\$1|IN \(.*\)).*`).
			WillReturnRows(joinRows)

		childRows := sqlmock.NewRows([]string{"id"})
		for id := range childIDs {
			childRows.AddRow(id)
		}
		m.ExpectQuery(`SELECT .* FROM "` + childTable + `".*WHERE .*"` + childTable + `"\."id" (=\s*\$1|IN \(.*\)).*`).
			WillReturnRows(childRows)
	}
}

/* ------------------Raising Status warning Helper------------------ */
func wantStatus(t *testing.T, got *http.Response, want int) {
	t.Helper()
	if got.StatusCode != want {
		t.Fatalf("status = %d, want %d; body=%s", got.StatusCode, want, string(readBody(t, got.Body)))
	}
}

/* ------------------Http setup ------------------ */

type httpInput struct {
	Method      string
	Path        string
	Body        []byte
	ContentType string
	UserID      *uint
}

func runHTTP(t *testing.T, app *fiber.App, in httpInput) *http.Response {
	t.Helper()

	var body *bytes.Reader
	if in.Body != nil {
		body = bytes.NewReader(in.Body)
	} else {
		body = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(in.Method, in.Path, body)
	if in.ContentType != "" {
		req.Header.Set("Content-Type", in.ContentType)
	}
	if in.UserID != nil {
		token := makeJWT(t, []byte(secretString), uint(*in.UserID))
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("app.Test: %v", err)
	}
	return resp
}

type fakeUploader struct {
	lastBucket   string
	lastFilename string
	lastData     []byte
	wantErr      error
}

func (f *fakeUploader) UploadBytes(_ context.Context, bucket, filename string, b []byte) (string, error) {
	f.lastBucket = bucket
	f.lastFilename = filename
	f.lastData = append([]byte(nil), b...)
	if f.wantErr != nil {
		return "", f.wantErr
	}
	return bucket + "/" + filename, nil
}

func newCtxWithUploader(u storage.Uploader) (*fiber.App, *fiber.Ctx) {
	app := fiber.New()
	rc := new(fasthttp.RequestCtx)
	ctx := app.AcquireCtx(rc)
	ctx.Locals("minio", u)
	return app, ctx
}

// 1×1 transparent PNG (raw base64, no data: prefix)
const tinyPNGRawBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMBAgK2V9sAAAAASUVORK5CYII="
