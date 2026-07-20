package middleware

import (
	"log"
	"time"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
	"github.com/gin-gonic/gin"
)

type logEntry struct {
	eventType  string
	username   string
	ip         string
	path       string
	method     string
	statusCode int
	details    string
	timestamp  time.Time
}

func writeAuditLog(entry logEntry) {
	_, err := db.DB.Exec(
		`INSERT INTO audit_logs (event_type, username, ip_address, request_path, request_method, status_code, details, timestamp)
						 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		entry.eventType,
		entry.username,
		entry.ip,
		entry.path,
		entry.method,
		entry.statusCode,
		entry.details,
		entry.timestamp,
	)
	if err != nil {
		log.Printf("[audit] failed to write log: %v", err)
	}
}

func WriteLoginAudit(eventType, username, ip, details string) {
	entry := logEntry{
		eventType:  eventType,
		username:   username,
		ip:         ip,
		path:       "/api/v1/auth/login",
		method:     "POST",
		statusCode: 0, // not meaningful for login-specific events
		details:    details,
		timestamp:  time.Now().UTC(),
	}
	go writeAuditLog(entry)
}

func AuditLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()

		username := "anonymous"
		if usr, exists := c.Get("username"); exists {
			if u, ok := usr.(string); ok && u != "" {
				username = u
			}
		}

		eventType := "API_ACCESS"
		if statusCode == 401 {
			eventType = "UNAUTHORIZED_ACCESS"
		} else if statusCode == 429 {
			eventType = "RATE_LIMITED"
		}

		entry := logEntry{
			eventType:  eventType,
			username:   username,
			ip:         ip,
			path:       path,
			method:     method,
			statusCode: statusCode,
			details:    "",
			timestamp:  time.Now().UTC(),
		}

		go writeAuditLog(entry)
	}
}
