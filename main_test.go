package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
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
