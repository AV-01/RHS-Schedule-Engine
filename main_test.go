package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
	"github.com/AV-01/RHS-Schedule-Engine/internal/handlers"
	"github.com/AV-01/RHS-Schedule-Engine/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// setupRouter configures a test router instance with standard middlewares
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(middleware.RateLimiter())

	// test endpoint that requires auth
	protected := r.Group("/api/test")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
	}
	return r
}

func TestMain(m *testing.M) {
	_ = godotenv.Load()
	os.Exit(m.Run())
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/test/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", w.Code)
	}
}

func TestAuthMiddleware_BadTokenFormat(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/test/ping", nil)
	req.Header.Set("Authorization", "InvalidFormatToken")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidDemoToken(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/test/ping", nil)
	req.Header.Set("Authorization", "Bearer demo-key")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	expectedBody := `{"message":"pong"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body %s, got %s", expectedBody, w.Body.String())
	}
}

func TestRateLimiter(t *testing.T) {
	middleware.ResetClients()
	router := setupRouter()

	// Make 61 requests to trigger the rate limiter (limit is 60)
	for i := 0; i < 65; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/test/ping", nil)
		req.Header.Set("Authorization", "Bearer demo-key")
		router.ServeHTTP(w, req)

		if i >= 60 {
			if w.Code != http.StatusTooManyRequests {
				t.Fatalf("Request %d: Expected rate limiter block (429), but got status %d", i+1, w.Code)
			}
		} else {
			if w.Code != http.StatusOK {
				t.Fatalf("Request %d: Expected status 200, got %d", i+1, w.Code)
			}
		}
	}
}

func TestDatabaseConnection(t *testing.T) {
	if os.Getenv("SUPABASE_DB_URL") == "" {
		t.Skip("Skipping DB test: SUPABASE_DB_URL not set")
	}

	db.Init()
	err := db.DB.Ping()
	if err != nil {
		t.Errorf("Database ping failed: %v", err)
	}
}

func setupAppRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/students", handlers.GetStudents)
		protected.GET("/students/:id", handlers.GetStudent)
		protected.GET("/students/:id/schedules", handlers.GetStudentSchedules)
		protected.GET("/classes", handlers.GetClasses)
		protected.GET("/teachers", handlers.GetTeachers)
		protected.GET("/teachers/:name/schedules", handlers.GetTeacherSchedule)
	}
	return r
}

func TestDemoStudentsEndpoint(t *testing.T) {
	router := setupAppRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/students?name=alex", nil)
	req.Header.Set("Authorization", "Bearer demo-key")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d. Body: %s", w.Code, w.Body.String())
	}

	body := w.Body.String()
	if !strings.Contains(body, "Alex") || !strings.Contains(body, "demo-student-001-uuid") {
		t.Errorf("Expected body to contain mock student Alex, got: %s", body)
	}
}

func TestDemoStudentDetailsAndSchedules(t *testing.T) {
	router := setupAppRouter()

	// GET Student Details
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/students/demo-student-001-uuid", nil)
	req1.Header.Set("Authorization", "Bearer demo-key")
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w1.Code)
	}
	if !strings.Contains(w1.Body.String(), "alex.morgan") {
		t.Errorf("Expected student username alex.morgan, got: %s", w1.Body.String())
	}

	// GET Student Schedules
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/students/demo-student-001-uuid/schedules", nil)
	req2.Header.Set("Authorization", "Bearer demo-key")
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w2.Code)
	}
	if !strings.Contains(w2.Body.String(), "AP Computer Science A") {
		t.Errorf("Expected AP Computer Science A in schedule, got: %s", w2.Body.String())
	}
}

func TestDemoClassesEndpoint(t *testing.T) {
	router := setupAppRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/classes", nil)
	req.Header.Set("Authorization", "Bearer demo-key")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "Algebra II") {
		t.Errorf("Expected class Algebra II in response, got: %s", w.Body.String())
	}
}

func TestDemoTeachersEndpoint(t *testing.T) {
	router := setupAppRouter()

	// GET Teachers List
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest(http.MethodGet, "/api/v1/teachers", nil)
	req1.Header.Set("Authorization", "Bearer demo-key")
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w1.Code)
	}
	if !strings.Contains(w1.Body.String(), "Stark, T") {
		t.Errorf("Expected teacher Stark, T in list, got: %s", w1.Body.String())
	}

	// GET Teacher Schedule
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest(http.MethodGet, "/api/v1/teachers/Stark,%20T/schedules", nil)
	req2.Header.Set("Authorization", "Bearer demo-key")
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", w2.Code)
	}
	if !strings.Contains(w2.Body.String(), "AP Computer Science A") {
		t.Errorf("Expected teacher class in schedule, got: %s", w2.Body.String())
	}
}
