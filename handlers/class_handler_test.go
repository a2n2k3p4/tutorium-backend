package handlers

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

/* ------------------ CreateClass ------------------ */

// 201
func TestCreateClass_OK(t *testing.T) {
	table := "classes"
	userID := uint(42)
	teacherID := uint(40)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, true)(mock)
			ExpInsertReturningID(table, 1)(mock)

			req := jsonBody(models.Class{
				TeacherID:        teacherID,
				ClassName:        "Testing",
				ClassDescription: "Lorem Ipsum",
				BannerPictureURL: "",
			})
			*payload = req
			*uID = userID
		},
		http.StatusCreated,
		http.MethodPost,
		"/classes/",
	)
}

// 400
func TestCreateClass_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		"/classes/",
	)
}

// 500
func TestCreateClass_DBError(t *testing.T) {
	table := "classes"
	userID := uint(42)
	teacherID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpInsertError(table, fmt.Errorf("db insert failed"))(mock)
			req := jsonBody(models.Class{
				TeacherID:        teacherID,
				ClassName:        "Testing",
				ClassDescription: "Lorem Ipsum",
				BannerPictureURL: "",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPost,
		"/classes/",
	)
}

/* ------------------ GetClasses ------------------ */

// 200
func TestGetClasses_OK(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListRows("classes", []string{"id"}, []any{1}, []any{2})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		"/classes/",
	)
}

// 500
func TestGetClasses_DBError(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpListError("classes", fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		"/classes/",
	)
}

/* ------------------ GetClass ------------------ */

// 200
func TestGetClass_OK(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(5)
	classCategoryID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDFound(table, classID, []string{"id"}, []any{classID})(mock)
			ExpPreloadTeacherForClassAny()(mock)
			ExpPreloadM2M("class_class_categories", "class_categories", "class_id", "class_category_id", [][2]any{{classID, classCategoryID}})(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 404
func TestGetClass_NotFound(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDEmpty(table, classID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 500
func TestGetClass_DBError(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(7)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpSelectByIDError(table, classID, fmt.Errorf("select failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodGet,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 400
func TestGetClass_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodGet,
		"/classes/not-an-int",
	)
}

/* ------------------ UpdateClass ------------------ */

// 200
func TestUpdateClass_OK(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(1)
	teacherID := uint(40)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, false, true, false)(mock)
			ExpSelectByIDFound(
				table,
				classID,
				[]string{"id", "teacher_id", "class_name", "class_description", "bannerPicture_url", "price", "rating"},
				[]any{classID, teacherID, "testing", "Lorem", "", 100, 4},
			)(mock)
			ExpPreloadField(table, []string{"id"}, []any{classID})(mock)
			ExpUpdateOK(table)(mock)

			req := jsonBody(models.Class{
				TeacherID:        teacherID,
				ClassName:        "Testing",
				ClassDescription: "Lorem Ipsum",
				BannerPictureURL: "",
			})
			*payload = req
			*uID = userID
		},
		http.StatusOK,
		http.MethodPut,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 404
func TestUpdateClass_NotFound(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, classID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodPut,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 500
func TestUpdateClass_DBError(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(1)
	teacherID := uint(40)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, classID, []string{"id"}, []any{classID})(mock)
			ExpUpdateError(table, fmt.Errorf("update failed"))(mock)

			req := jsonBody(models.Class{
				TeacherID:        teacherID,
				ClassName:        "Testing",
				ClassDescription: "Lorem Ipsum",
				BannerPictureURL: "",
			})
			*payload = req
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodPut,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 400
func TestUpdateClass_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodPut,
		"/classes/not-an-int",
	)
}

/* ------------------ DeleteClass ------------------ */

// 200 (soft delete)
func TestDeleteClass_OK_SoftDelete(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, classID, []string{"id"}, []any{classID})(mock)
			ExpClearCategoriesForClass(classID)(mock)
			ExpSoftDeleteOK(table)(mock)
			*uID = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 404
func TestDeleteClass_NotFound(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(12345)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty(table, classID)(mock)
			*uID = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 500
func TestDeleteClass_DBError(t *testing.T) {
	table := "classes"
	userID := uint(42)
	classID := uint(5)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound(table, classID, []string{"id"}, []any{classID})(mock)
			ExpClearCategoriesForClass(classID)(mock)
			ExpSoftDeleteError(table, fmt.Errorf("update failed"))(mock)
			*uID = userID
		},
		http.StatusInternalServerError,
		http.MethodDelete,
		fmt.Sprintf("/classes/%d", classID),
	)
}

// 400
func TestDeleteClass_BadRequest(t *testing.T) {
	userID := uint(42)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uID *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			*uID = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		"/classes/not-an-int",
	)
}

// ---------- AddClassCategories ----------
func TestAddClassCategories_OK(t *testing.T) {
	userID := uint(42)
	classID := uint(10)
	// send three category IDs in body
	newCats := []uint{2, 4, 6}

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound("classes", classID, []string{"id"}, []any{classID})(mock)
			ExpAssociationFindClassCategoriesEmpty(classID)(mock)
			ExpSelectCategoriesByIDs(newCats...)(mock)
			ExpAppendClassCategories(classID, newCats...)(mock)
			ExpPreloadClassWithCategories(classID, newCats...)(mock)
			type body struct {
				CategoryIDs []int `json:"class_category_ids"`
			}
			b, _ := json.Marshal(body{CategoryIDs: []int{2, 4, 6}})
			*payload = b
			*uid = userID
		},
		http.StatusOK,
		http.MethodPost,
		fmt.Sprintf("/classes/%d/categories", classID),
	)
}

func TestAddClassCategories_NoIDs_BadRequest(t *testing.T) {
	userID := uint(42)
	classID := uint(10)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			type body struct {
				CategoryIDs []int `json:"class_category_ids"`
			}
			b, _ := json.Marshal(body{CategoryIDs: []int{}})
			*payload = b
			*uid = userID
		},
		http.StatusBadRequest,
		http.MethodPost,
		fmt.Sprintf("/classes/%d/categories", classID),
	)
}

func TestAddClassCategories_NotFound(t *testing.T) {
	userID := uint(42)
	classID := uint(999)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty("classes", classID)(mock)
			type body struct {
				CategoryIDs []int `json:"class_category_ids"`
			}
			b, _ := json.Marshal(body{CategoryIDs: []int{2}})
			*payload = b
			*uid = userID
		},
		http.StatusNotFound,
		http.MethodPost,
		fmt.Sprintf("/classes/%d/categories", classID),
	)
}

// ---------- DeleteClassCategories ----------

func TestDeleteClassCategories_OK(t *testing.T) {
	userID := uint(42)
	classID := uint(10)
	delCats := []uint{2, 4}

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDFound("classes", classID, []string{"id"}, []any{classID})(mock)
			ExpSelectCategoriesByIDs(delCats...)(mock)
			ExpDeleteClassCategories(classID, delCats...)(mock)
			ExpPreloadClassWithCategories(classID /*none*/)(mock)
			type body struct {
				CategoryIDs []int `json:"class_category_ids"`
			}
			b, _ := json.Marshal(body{CategoryIDs: []int{2, 4}})
			*payload = b
			*uid = userID
		},
		http.StatusOK,
		http.MethodDelete,
		fmt.Sprintf("/classes/%d/categories", classID),
	)
}

func TestDeleteClassCategories_NoIDs_BadRequest(t *testing.T) {
	userID := uint(42)
	classID := uint(10)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			type body struct {
				CategoryIDs []int `json:"class_category_ids"`
			}
			b, _ := json.Marshal(body{CategoryIDs: []int{}})
			*payload = b
			*uid = userID
		},
		http.StatusBadRequest,
		http.MethodDelete,
		fmt.Sprintf("/classes/%d/categories", classID),
	)
}

func TestDeleteClassCategories_NotFound(t *testing.T) {
	userID := uint(42)
	classID := uint(999)
	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, true, false, false)(mock)
			ExpSelectByIDEmpty("classes", classID)(mock)
			type body struct {
				CategoryIDs []int `json:"class_category_ids"`
			}
			b, _ := json.Marshal(body{CategoryIDs: []int{2}})
			*payload = b
			*uid = userID
		},
		http.StatusNotFound,
		http.MethodDelete,
		fmt.Sprintf("/classes/%d/categories", classID),
	)
}

// ---------- GetClassCategoriesByClassID ----------
func TestGetClassCategoriesByClassID_OK(t *testing.T) {
	userID := uint(42)
	classID := uint(10)
	cats := []uint{1, 3, 5}

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			ExpPreloadClassCategoriesOrdered(classID, cats...)(mock)

			*uid = userID
		},
		http.StatusOK,
		http.MethodGet,
		fmt.Sprintf("/classes/%d/categories", classID),
	)
}

func TestGetClassCategoriesByClassID_NotFound(t *testing.T) {
	userID := uint(42)
	classID := uint(999)

	RunInDifferentStatus(t,
		func(t *testing.T, mock sqlmock.Sqlmock, gdb *gorm.DB, app *fiber.App, payload *[]byte, uid *uint) {
			ExpAuthUser(userID, false, false, false)(mock)
			mock.ExpectQuery(q(`SELECT * FROM "classes" WHERE "classes"."id" = $1 AND "classes"."deleted_at" IS NULL ORDER BY "classes"."id" LIMIT 1`)).
				WithArgs(driver.Value(int64(classID))).
				WillReturnRows(sqlmock.NewRows([]string{"id"}))
			*uid = userID
		},
		http.StatusNotFound,
		http.MethodGet,
		fmt.Sprintf("/classes/%d/categories", classID),
	)
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
