package handlers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/stretchr/testify/require"
)

/* ------------------ CreateClass ------------------ */

// 201
func TestCreateClass_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		teacherID := uint(40)

		ExpAuthUser(userID, false, false, true)(mock)
		ExpInsertReturningID(table, 1)(mock)

		app := setupApp(gdb)

		payload := models.Class{
			TeacherID:        teacherID,
			ClassName:        "Testing",
			ClassDescription: "Lorem Ipsum",
			BannerPictureURL: "",
			Price:            50,
			Rating:           5,
		}

		resp := runHTTP(t, app, httpInput{
			Method:      http.MethodPost,
			Path:        "/classes/",
			Body:        jsonBody(payload),
			ContentType: "application/json",
			UserID:      &userID,
		})
		wantStatus(t, resp, http.StatusCreated)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestCreateClass_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/classes/",
			Body: []byte(`{invalid-json}`), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestCreateClass_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPost, Path: "/classes/",
			Body: []byte(`{}`), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetClasses ------------------ */
// 200
func TestGetClasses_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpListRows("classes", []string{"id"}, []any{1}, []any{2})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/classes/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetClasses_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpListError("classes", fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/classes/", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ GetClass ------------------ */
// 200
func TestGetClass_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		classCategoryID := uint(7)
		classID := uint(5)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
		ExpPreloadM2M("class_class_categories", "class_categories", "class_id", "class_category_id", [][2]any{{classCategoryID, classID}})(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/classes/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestGetClass_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		classCategoryID := uint(999)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDEmpty(table, classCategoryID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/classes/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestGetClass_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		classCategoryID := uint(7)

		ExpAuthUser(userID, false, false, false)(mock)
		ExpSelectByIDError(table, classCategoryID, fmt.Errorf("select failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: fmt.Sprintf("/classes/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestGetClass_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, false, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodGet, Path: "/classes/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ UpdateClass ------------------ */

// 200
func TestUpdateClass_OK(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"

		userID := uint(42)
		classID := uint(1)
		teacherID := uint(40)

		ExpAuthUser(userID, false, true, false)(mock)
		ExpSelectByIDFound(
			table,
			classID,
			[]string{"id", "teacher_id", "class_name", "class_description", "bannerPicture_url", "price", "rating"},
			[]any{classID, teacherID, "testing", "Lorem", "", 100, 4},
		)(mock)
		ExpUpdateOK(table)(mock)
		ExpPreloadField(table, []string{"id"}, []any{classID})(mock)

		app := setupApp(gdb)
		payload := models.Class{
			TeacherID:        teacherID,
			ClassName:        "Testing",
			ClassDescription: "Lorem Ipsum",
			BannerPictureURL: "",
			Price:            50,
			Rating:           5,
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/classes/%d", classID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestUpdateClass_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		classCategoryID := uint(12345)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDEmpty(table, classCategoryID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/classes/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestUpdateClass_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		classCategoryID := uint(1)
		teacherID := uint(40)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id", "class_category"}, []any{classCategoryID, "test"})(mock)
		ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		payload := models.Class{
			TeacherID:        teacherID,
			ClassName:        "Testing",
			ClassDescription: "Lorem Ipsum",
			BannerPictureURL: "",
			Price:            50,
			Rating:           5,
		}
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: fmt.Sprintf("/classes/%d", classCategoryID),
			Body: jsonBody(payload), ContentType: "application/json", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestUpdateClass_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodPut, Path: "/classes/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ DeleteClass ------------------ */

// 200 (soft delete)
func TestDeleteClass_OK_SoftDelete(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		classCategoryID := uint(5)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
		ExpSoftDeleteOK(table)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/classes/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusOK)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 404
func TestDeleteClass_NotFound(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		classCategoryID := uint(12345)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDEmpty(table, classCategoryID)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/classes/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusNotFound)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 500
func TestDeleteClass_DBError(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)

		table := "classes"
		userID := uint(42)
		classCategoryID := uint(5)

		ExpAuthUser(userID, true, false, false)(mock)
		ExpSelectByIDFound(table, classCategoryID, []string{"id"}, []any{classCategoryID})(mock)
		ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: fmt.Sprintf("/classes/%d", classCategoryID), UserID: &userID,
		})
		wantStatus(t, resp, http.StatusInternalServerError)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

// 400
func TestDeleteClass_BadRequest(t *testing.T) {

	RunInDifferentStatus(t, func(t *testing.T) {
		mock, gdb, cleanup := setupMockGorm(t)
		defer cleanup()
		mock.MatchExpectationsInOrder(false)
		userID := uint(42)

		ExpAuthUser(userID, true, false, false)(mock)

		app := setupApp(gdb)
		resp := runHTTP(t, app, httpInput{
			Method: http.MethodDelete, Path: "/classes/not-an-int", UserID: &userID,
		})
		wantStatus(t, resp, http.StatusBadRequest)
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Fatalf("unmet expectations: %v", err)
		}
	})
}

/* ------------------ processBannerPicture ------------------ */

func TestProcessBannerPicture_Success(t *testing.T) {
	fu := &fakeUploader{}
	app, ctx := newCtxWithUploader(fu)
	defer app.ReleaseCtx(ctx)

	u := &models.Class{BannerPictureURL: tinyPNGRawBase64}

	err := processBannerPicture(ctx, u)
	require.NoError(t, err)

	require.Equal(t, "classes", fu.lastBucket)
	require.NotEmpty(t, fu.lastFilename)
	require.NotEmpty(t, fu.lastData)

	require.Regexp(t, regexp.MustCompile(`^classes/\d+\.(png|jpg|jpeg|gif|webp|bin)$`), u.BannerPictureURL)
}

func TestProcessBannerPicture_SkipEmptyAndHTTP(t *testing.T) {
	fu := &fakeUploader{}
	app, ctx := newCtxWithUploader(fu)
	defer app.ReleaseCtx(ctx)

	u1 := &models.Class{BannerPictureURL: ""}
	require.NoError(t, processBannerPicture(ctx, u1))
	require.Equal(t, "", u1.BannerPictureURL)
	require.Equal(t, "", fu.lastFilename)

	u2 := &models.Class{BannerPictureURL: "http://example.com/pic.png"}
	require.NoError(t, processBannerPicture(ctx, u2))
	require.Equal(t, "http://example.com/pic.png", u2.BannerPictureURL)
	require.Equal(t, "", fu.lastFilename)
}

func TestProcessBannerPicture_InvalidBase64(t *testing.T) {
	fu := &fakeUploader{}
	app, ctx := newCtxWithUploader(fu)
	defer app.ReleaseCtx(ctx)

	u := &models.Class{BannerPictureURL: "!!!not-b64!!!"}
	err := processBannerPicture(ctx, u)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid base64 image")
}

func TestProcessBannerPicture_InvalidImage(t *testing.T) {
	fu := &fakeUploader{}
	app, ctx := newCtxWithUploader(fu)
	defer app.ReleaseCtx(ctx)

	s := base64.StdEncoding.EncodeToString([]byte("hello"))
	u := &models.Class{BannerPictureURL: s}
	err := processBannerPicture(ctx, u)
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid image")
}
