package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/a2n2k3p4/tutorium-backend/config"
	"github.com/a2n2k3p4/tutorium-backend/middlewares"
	"github.com/a2n2k3p4/tutorium-backend/models"
	"github.com/gofiber/fiber/v2"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	integApp      *fiber.App
	integDB       *gorm.DB
	pgC           tc.Container
	uniqueCounter atomic.Uint64
)

/* ------------------ Main ------------------ */
func TestMain(m *testing.M) {
	_ = os.Setenv("STATUS", "development")

	ctx := context.Background()
	req := tc.ContainerRequest{
		Image:        "postgres:17",
		Env:          map[string]string{"POSTGRES_PASSWORD": "password", "POSTGRES_USER": "user", "POSTGRES_DB": "tutorium"},
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(90 * time.Second),
	}

	var err error
	pgC, err = tc.GenericContainer(ctx, tc.GenericContainerRequest{ContainerRequest: req, Started: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start postgres container: %v\n", err)
		os.Exit(1)
	}

	host, err := pgC.Host(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "container host error: %v\n", err)
		os.Exit(1)
	}
	port, err := pgC.MappedPort(ctx, "5432/tcp")
	if err != nil {
		fmt.Fprintf(os.Stderr, "container port map error: %v\n", err)
		os.Exit(1)
	}

	cfg := &config.Config{
		DBUser:     "user",
		DBPassword: "password",
		DBHost:     host,
		DBPort:     string(port.Port()),
		DBName:     "tutorium",
	}
	integDB, err = connectDB_with_silent(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect db error: %v\n", err)
		os.Exit(1)
	}

	models.Migrate(integDB)
	middlewares.Status = func() string { return "development" }

	integApp = fiber.New()
	integApp.Use(middlewares.DBMiddleware(integDB))
	integApp.Use(middlewares.MinioMiddleware(dummyUploader{}))
	integApp.Use(func(c *fiber.Ctx) error {
		c.Locals("omise", nil)
		return c.Next()
	})
	AllRoutes(integApp)

	code := m.Run()

	if pgC != nil {
		if err := pgC.Terminate(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "terminate container error: %v\n", err)
		}
	}

	os.Exit(code)
}

/* ------------------ Connecting and middleware Helper ------------------ */
type dummyUploader struct{}

func (dummyUploader) UploadBytes(ctx context.Context, folder, filename string, b []byte) (string, error) {
	return fmt.Sprintf("stub://%s/%s", folder, filename), nil
}

func connectDB_with_silent(cfg *config.Config) (*gorm.DB, error) {
	dbUrl := cfg.DBUrl()

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil
}

/* ------------------ API Request Helper ------------------ */
func newJSONRequest(t *testing.T, method, target string, payload any) *http.Request {
	t.Helper()

	var body io.Reader
	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		body = bytes.NewReader(data)
	}
	req := httptest.NewRequest(method, target, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func performRequest(t *testing.T, req *http.Request) *http.Response {
	t.Helper()

	resp, err := integApp.Test(req, -1)
	if err != nil {
		t.Fatalf("fiber test error: %v", err)
	}
	return resp
}

func requireStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Fatalf("unexpected status %d (want %d): %s", resp.StatusCode, want, string(body))
	}
}

func decodeJSON(t *testing.T, resp *http.Response, out any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func jsonRequestExpect(t *testing.T, method, target string, payload any, wantStatus int, out any) {
	t.Helper()

	resp := performRequest(t, newJSONRequest(t, method, target, payload))
	requireStatus(t, resp, wantStatus)

	if out == nil {
		resp.Body.Close()
		return
	}

	decodeJSON(t, resp, out)
}

/* ------------------ CRUD Test Helper ------------------ */
func createJSONResource[T any](t *testing.T, target string, payload any, wantStatus int) T {
	t.Helper()

	var out T
	jsonRequestExpect(t, http.MethodPost, target, payload, wantStatus, &out)
	return out
}

func getJSONResource[T any](t *testing.T, target string, wantStatus int) T {
	t.Helper()

	var out T
	jsonRequestExpect(t, http.MethodGet, target, nil, wantStatus, &out)
	return out
}

func updateJSONResource(t *testing.T, target string, payload any, wantStatus int) {
	t.Helper()
	jsonRequestExpect(t, http.MethodPut, target, payload, wantStatus, nil)
}

func deleteJSONResource(t *testing.T, target string, wantStatus int) {
	t.Helper()
	jsonRequestExpect(t, http.MethodDelete, target, nil, wantStatus, nil)
}

func nextSequence() uint64 {
	return uniqueCounter.Add(1)
}

func uniqueSuffix() string {
	return fmt.Sprintf("%06d", nextSequence())
}

func randomStudentID() string {
	return fmt.Sprintf("b67%08d", nextSequence())
}

func randomPhoneNumber() string {
	return fmt.Sprintf("+669%08d", nextSequence())
}

func randomEmail(prefix string) string {
	return fmt.Sprintf("%s-%d@example.com", prefix, nextSequence())
}

func createTestUser(t *testing.T) models.User {
	t.Helper()

	payload := map[string]any{
		"student_id":   randomStudentID(),
		"first_name":   "Integration",
		"last_name":    "User",
		"gender":       "Other",
		"phone_number": randomPhoneNumber(),
	}

	return createJSONResource[models.User](t, "/users/", payload, http.StatusCreated)
}

func createTestTeacher(t *testing.T) (models.Teacher, models.User) {
	t.Helper()

	user := createTestUser(t)
	payload := map[string]any{
		"user_id":     user.ID,
		"description": "Integration teacher",
		"email":       randomEmail("teacher"),
	}
	teacher := createJSONResource[models.Teacher](t, "/teachers/", payload, http.StatusCreated)
	return teacher, user
}

func createTestClass(t *testing.T, teacherID uint) models.Class {
	t.Helper()

	payload := map[string]any{
		"teacher_id":        teacherID,
		"class_name":        fmt.Sprintf("Integration Class %s", uniqueSuffix()),
		"class_description": "Integration class description",
		"rating":            4.5,
		"categories": []map[string]string{
			{"class_category": "Mathematics"},
		},
	}
	return createJSONResource[models.Class](t, "/classes/", payload, http.StatusCreated)
}

func createTestLearner(t *testing.T) (models.Learner, models.User) {
	t.Helper()

	user := createTestUser(t)
	if user.Learner == nil {
		t.Fatalf("expected learner to be created automatically for user %d", user.ID)
	}
	return *user.Learner, user
}

func createTestClassSession(t *testing.T, classID uint) models.ClassSession {
	t.Helper()
	now := time.Now()
	payload := map[string]any{
		"class_id":            classID,
		"description":         fmt.Sprintf("Session %s", uniqueSuffix()),
		"price":               1234,
		"learner_limit":       30,
		"enrollment_deadline": now.Add(48 * time.Hour).Format(time.RFC3339Nano),
		"class_start":         now.Add(72 * time.Hour).Format(time.RFC3339Nano),
		"class_finish":        now.Add(96 * time.Hour).Format(time.RFC3339Nano),
		"class_status":        "upcoming",
	}
	return createJSONResource[models.ClassSession](t, "/class_sessions/", payload, http.StatusCreated)
}

func createTestEnrollment(t *testing.T, learnerID, classSessionID uint) models.Enrollment {
	t.Helper()
	payload := map[string]any{
		"learner_id":        learnerID,
		"class_session_id":  classSessionID,
		"enrollment_status": "active",
	}
	return createJSONResource[models.Enrollment](t, "/enrollments/", payload, http.StatusCreated)
}

func createTestNotification(t *testing.T, userID uint) models.Notification {
	t.Helper()
	payload := map[string]any{
		"user_id":                  userID,
		"notification_type":        "integration",
		"notification_description": "Integration notification",
		"notification_date":        time.Now().Format(time.RFC3339Nano),
		"read_flag":                false,
	}
	return createJSONResource[models.Notification](t, "/notifications/", payload, http.StatusCreated)
}

func createTestReport(t *testing.T, reporterID, reportedID, classSessionID uint) models.Report {
	t.Helper()
	payload := map[string]any{
		"report_user_id":     reporterID,
		"reported_user_id":   reportedID,
		"class_session_id":   classSessionID,
		"report_type":        "behavior",
		"report_reason":      "integration",
		"report_description": "integration report description",
		"report_date":        time.Now().Format(time.RFC3339Nano),
		"report_status":      "pending",
	}
	return createJSONResource[models.Report](t, "/reports/", payload, http.StatusCreated)
}

func createTestReview(t *testing.T, learnerID, classID uint) models.Review {
	t.Helper()
	payload := map[string]any{
		"learner_id": learnerID,
		"class_id":   classID,
		"rating":     5,
		"comment":    "Integration review",
	}
	return createJSONResource[models.Review](t, "/reviews/", payload, http.StatusCreated)
}

func createTestBanLearner(t *testing.T, learnerID uint) models.BanDetailsLearner {
	t.Helper()
	now := time.Now()
	payload := map[string]any{
		"learner_id":      learnerID,
		"ban_start":       now.Add(-1 * time.Hour).Format(time.RFC3339Nano),
		"ban_end":         now.Add(24 * time.Hour).Format(time.RFC3339Nano),
		"ban_description": "integration ban learner",
	}
	return createJSONResource[models.BanDetailsLearner](t, "/banlearners/", payload, http.StatusCreated)
}

func createTestBanTeacher(t *testing.T, teacherID uint) models.BanDetailsTeacher {
	t.Helper()
	now := time.Now()
	payload := map[string]any{
		"teacher_id":      teacherID,
		"ban_start":       now.Add(-1 * time.Hour).Format(time.RFC3339Nano),
		"ban_end":         now.Add(24 * time.Hour).Format(time.RFC3339Nano),
		"ban_description": "integration ban teacher",
	}
	return createJSONResource[models.BanDetailsTeacher](t, "/banteachers/", payload, http.StatusCreated)
}
