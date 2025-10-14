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
	"regexp"
	"strings"
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

/* ------------------ Tx Helpers ------------------ */
func ExpTxBegin() func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) { mock.ExpectBegin() }
}
func ExpTxCommit() func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) { mock.ExpectCommit() }
}
func ExpTxRollback() func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) { mock.ExpectRollback() }
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

/* ------------------ SQL helper ------------------ */
func q(s string) string { return regexp.QuoteMeta(s) }

func toVals(ids ...uint) []driver.Value {
	out := make([]driver.Value, len(ids))
	for i, id := range ids {
		out[i] = driver.Value(int64(id))
	}
	return out
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

func ExpPreloadSessionsForClassOrdered(classID uint, cols ...string) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		if len(cols) == 0 {
			cols = []string{"id", "class_id"}
		}
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "class_sessions" WHERE class_id = $1 AND "class_sessions"."deleted_at" IS NULL ORDER BY class_start ASC`,
		)).
			WithArgs(driver.Value(int64(classID))).
			WillReturnRows(sqlmock.NewRows(cols))
	}
}

func ExpClearJoinByFK(joinTable, fk string, id uint, rows int64) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(
			fmt.Sprintf(`DELETE FROM "%s" WHERE "%s"."%s" = $1`, joinTable, joinTable, fk),
		)).
			WithArgs(driver.Value(int64(id))).
			WillReturnResult(sqlmock.NewResult(0, rows))
	}
}

func ExpSelectLearnerByUserIDEmpty(userID uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "learners" WHERE user_id = $1 AND "learners"."deleted_at" IS NULL ORDER BY "learners"."id" LIMIT 1`,
		)).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}
}

func ExpClearInterestedForLearner(learnerID uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM "interested_class_categories" WHERE "interested_class_categories"."learner_id" = $1`,
		)).
			WithArgs(learnerID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func ExpClearCategoriesForClass(classID uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM "class_class_categories" WHERE "class_class_categories"."class_id" = $1`,
		)).
			WithArgs(classID).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func ExpPreloadLearnerInterestedEmpty(learnerID uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(q(`SELECT * FROM "interested_class_categories" WHERE "interested_class_categories"."learner_id" = $1`)).
			WithArgs(driver.Value(int64(learnerID))).
			WillReturnRows(sqlmock.NewRows([]string{"learner_id", "class_category_id"}))
	}
}

func ExpFetchAllActiveEmpty() func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`FROM\s+"classes"\s+JOIN\s+class_sessions\s+cs`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}
}

func ExpFetchAllActiveWithClasses(pairs [][2]uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		r := sqlmock.NewRows([]string{"id", "teacher_id"})
		classIDs := make([]uint, 0, len(pairs))
		teacherIDs := make([]uint, 0, len(pairs))
		for _, p := range pairs {
			cid, tid := p[0], p[1]
			classIDs = append(classIDs, cid)
			teacherIDs = append(teacherIDs, tid)
			r.AddRow(int64(cid), int64(tid))
		}
		mock.ExpectQuery(`SELECT\s+.*classes\.\*\s+FROM\s+"classes"\s+JOIN\s+class_sessions\s+cs.*`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(r)
		ph := make([]string, len(teacherIDs))
		for i := range teacherIDs {
			ph[i] = fmt.Sprintf(`\$%d`, i+1)
		}
		mock.ExpectQuery(`SELECT \* FROM "teachers" WHERE "teachers"\."id" IN \(` + strings.Join(ph, ",") + `\).*`).
			WithArgs(toVals(teacherIDs...)...).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}))
		ph = make([]string, len(classIDs))
		for i := range classIDs {
			ph[i] = fmt.Sprintf(`\$%d`, i+1)
		}
		mock.ExpectQuery(`SELECT \* FROM "class_class_categories" WHERE "class_class_categories"\."class_id" IN \(` + strings.Join(ph, ",") + `\)`).
			WithArgs(toVals(classIDs...)...).
			WillReturnRows(sqlmock.NewRows([]string{"class_id", "class_category_id"}))
	}
}

func ExpRecommendedEmpty() func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(`FROM\s+"classes"\s+JOIN\s+class_sessions\s+cs.*JOIN\s+class_class_categories\s+ccc`).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))
	}
}

func ExpSelectCategoriesByIDs(catIDs ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		ph := make([]string, len(catIDs))
		for i := range catIDs {
			ph[i] = fmt.Sprintf(`\$%d`, i+1)
		}
		rows := sqlmock.NewRows([]string{"id", "class_category"})
		for _, cid := range catIDs {
			rows.AddRow(int64(cid), fmt.Sprintf("Cat-%d", cid))
		}
		mock.ExpectQuery(q(fmt.Sprintf(
			`SELECT * FROM "class_categories" WHERE id IN (%s)`, strings.Join(ph, ","),
		))).
			WithArgs(toVals(catIDs...)...).
			WillReturnRows(rows)
	}
}

func ExpAppendLearnerInterests(learnerID uint, catIDs ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(q(`INSERT INTO "interested_class_categories" ("learner_id","class_category_id")`)).
			WillReturnResult(sqlmock.NewResult(0, int64(len(catIDs))))
	}
}

func ExpDeleteLearnerInterests(learnerID uint, catIDs ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		ph := make([]string, len(catIDs))
		for i := range catIDs {
			ph[i] = fmt.Sprintf(`\$%d`, i+2)
		}
		args := append([]driver.Value{driver.Value(int64(learnerID))}, toVals(catIDs...)...)
		mock.ExpectExec(
			`DELETE FROM "interested_class_categories" WHERE "interested_class_categories"\."learner_id" = \$1 AND "interested_class_categories"\."class_category_id" IN \(` + strings.Join(ph, ",") + `\)`,
		).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, int64(len(catIDs))))
	}
}

func ExpPreloadLearnerInterestedOrdered(learnerID uint, catIDs ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(q(`SELECT * FROM "interested_class_categories" WHERE "interested_class_categories"."learner_id" = $1`)).
			WithArgs(driver.Value(int64(learnerID))).
			WillReturnRows(func() *sqlmock.Rows {
				r := sqlmock.NewRows([]string{"learner_id", "class_category_id"})
				for _, cid := range catIDs {
					r.AddRow(int64(learnerID), int64(cid))
				}
				return r
			}())
		if len(catIDs) > 0 {
			ph := make([]string, len(catIDs))
			for i := range catIDs {
				ph[i] = fmt.Sprintf(`\$%d`, i+1)
			}
			rows := sqlmock.NewRows([]string{"id", "class_category"})
			for _, cid := range catIDs {
				rows.AddRow(int64(cid), fmt.Sprintf("Cat-%d", cid))
			}
			mock.ExpectQuery(
				`SELECT \* FROM "class_categories" WHERE "class_categories"\."id" IN \(` + strings.Join(ph, ",") + `\)`,
			).
				WithArgs(toVals(catIDs...)...).
				WillReturnRows(rows)
		}
	}
}

func ExpPreloadTeacherForClassAny() func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(regexp.QuoteMeta(
			`SELECT * FROM "teachers" WHERE "teachers"."id" IN ($1) AND "teachers"."deleted_at" IS NULL`,
		)).
			WithArgs(sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id"}).AddRow(1, 1))
	}
}

func ExpPreloadLearnersInterestedEmpty(ids ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		if len(ids) == 0 {
			return
		}
		ph := make([]string, len(ids))
		args := make([]driver.Value, len(ids))
		for i, id := range ids {
			ph[i] = fmt.Sprintf("$%d", i+1)
			args[i] = driver.Value(int64(id))
		}
		sql := fmt.Sprintf(
			`SELECT * FROM "interested_class_categories" WHERE "interested_class_categories"."learner_id" IN (%s)`,
			strings.Join(ph, ","),
		)
		mock.ExpectQuery(regexp.QuoteMeta(sql)).
			WithArgs(args...).
			WillReturnRows(sqlmock.NewRows([]string{"learner_id", "class_category_id"}))
	}
}

func ExpClearClassesForCategory(classCategoryID uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM "class_class_categories" WHERE "class_class_categories"."class_category_id" = $1`,
		)).
			WithArgs(driver.Value(int64(classCategoryID))).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func ExpClearLearnersForCategory(classCategoryID uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(
			`DELETE FROM "interested_class_categories" WHERE "interested_class_categories"."class_category_id" = $1`,
		)).
			WithArgs(driver.Value(int64(classCategoryID))).
			WillReturnResult(sqlmock.NewResult(0, 1))
	}
}

func ExpAssociationFindClassCategoriesEmpty(classID uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(q(`SELECT * FROM "class_class_categories" WHERE "class_class_categories"."class_id" = $1`)).
			WithArgs(driver.Value(int64(classID))).
			WillReturnRows(sqlmock.NewRows([]string{"class_id", "class_category_id"}))
	}
}

func ExpAppendClassCategories(classID uint, catIDs ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectExec(q(`INSERT INTO "class_class_categories" ("class_id","class_category_id")`)).
			WillReturnResult(sqlmock.NewResult(0, int64(len(catIDs))))
	}
}

func ExpPreloadClassWithCategories(classID uint, catIDs ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(q(`SELECT * FROM "classes" WHERE "classes"."id" = $1 AND "classes"."deleted_at" IS NULL ORDER BY "classes"."id" LIMIT 1`)).
			WithArgs(driver.Value(int64(classID))).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(classID)))

		mock.ExpectQuery(q(`SELECT * FROM "class_class_categories" WHERE "class_class_categories"."class_id" = $1`)).
			WithArgs(driver.Value(int64(classID))).
			WillReturnRows(func() *sqlmock.Rows {
				r := sqlmock.NewRows([]string{"class_id", "class_category_id"})
				for _, cid := range catIDs {
					r.AddRow(int64(classID), int64(cid))
				}
				return r
			}())
		if len(catIDs) > 0 {
			ph := make([]string, len(catIDs))
			for i := range catIDs {
				ph[i] = fmt.Sprintf(`\$%d`, i+1)
			}
			rows := sqlmock.NewRows([]string{"id", "class_category"})
			for _, cid := range catIDs {
				rows.AddRow(int64(cid), fmt.Sprintf("Cat-%d", cid))
			}
			mock.ExpectQuery(
				`SELECT \* FROM "class_categories" WHERE "class_categories"\."id" IN \(` + strings.Join(ph, ",") + `\)`,
			).
				WithArgs(toVals(catIDs...)...).
				WillReturnRows(rows)
		}
	}
}

func ExpPreloadClassCategoriesOrdered(classID uint, catIDs ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		mock.ExpectQuery(q(`SELECT * FROM "classes" WHERE "classes"."id" = $1 AND "classes"."deleted_at" IS NULL ORDER BY "classes"."id" LIMIT 1`)).
			WithArgs(driver.Value(int64(classID))).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(classID)))

		mock.ExpectQuery(q(`SELECT * FROM "class_class_categories" WHERE "class_class_categories"."class_id" = $1`)).
			WithArgs(driver.Value(int64(classID))).
			WillReturnRows(func() *sqlmock.Rows {
				r := sqlmock.NewRows([]string{"class_id", "class_category_id"})
				for _, cid := range catIDs {
					r.AddRow(int64(classID), int64(cid))
				}
				return r
			}())
		if len(catIDs) > 0 {
			ph := make([]string, len(catIDs))
			for i := range catIDs {
				ph[i] = fmt.Sprintf(`\$%d`, i+1)
			}
			rows := sqlmock.NewRows([]string{"id", "class_category"})
			for _, cid := range catIDs {
				rows.AddRow(int64(cid), fmt.Sprintf("Cat-%d", cid))
			}
			mock.ExpectQuery(
				`SELECT (id|.+) FROM "class_categories" WHERE "class_categories"\."id" IN \(` + strings.Join(ph, ",") + `\).*ORDER BY class_categories\.class_category`,
			).
				WithArgs(toVals(catIDs...)...).
				WillReturnRows(rows)
		}
	}
}

func ExpDeleteClassCategories(classID uint, catIDs ...uint) func(sqlmock.Sqlmock) {
	return func(mock sqlmock.Sqlmock) {
		ph := make([]string, len(catIDs))
		for i := range catIDs {
			ph[i] = fmt.Sprintf(`\$%d`, i+2)
		}
		args := append([]driver.Value{driver.Value(int64(classID))}, toVals(catIDs...)...)
		mock.ExpectExec(
			`DELETE FROM "class_class_categories" WHERE "class_class_categories"\."class_id" = \$1 AND "class_class_categories"\."class_category_id" IN \(` + strings.Join(ph, ",") + `\)`,
		).
			WithArgs(args...).
			WillReturnResult(sqlmock.NewResult(0, int64(len(catIDs))))
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
