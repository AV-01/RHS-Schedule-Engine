package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
	"github.com/AV-01/RHS-Schedule-Engine/internal/handlers"
	"github.com/AV-01/RHS-Schedule-Engine/internal/middleware"
)

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env found")
	}

	db.Init()

	router := gin.Default()

	router.Use(middleware.RateLimiter())

	public := router.Group("/api/v1")
	{
		public.POST("/auth/login", handlers.Login)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/students", handlers.GetStudents)
		protected.GET("/students/:id", handlers.GetStudent)
		protected.GET("/students/:id/schedules", handlers.GetStudentSchedules)

		protected.GET("/classes", handlers.GetClasses)

		protected.GET("/teachers", handlers.GetTeachers)
		protected.GET("/teachers/:name/schedules", handlers.GetTeacherSchedule)

	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
