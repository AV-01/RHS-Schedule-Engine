package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
	"github.com/gin-gonic/gin"
)

type AuditLog struct {
	ID            int    `json:"id"`
	EventType     string `json:"event_type"`
	Username      string `json:"username"`
	IPAddress     string `json:"ip_address"`
	RequestPath   string `json:"request_path"`
	RequestMethod string `json:"request_method"`
	StatusCode    int    `json:"status_code"`
	Timestamp     string `json:"timestamp"`
	Details       string `json:"details,omitempty"`
}

type PaginatedAuditLogs struct {
	Data  []AuditLog `json:"data"`
	Page  int        `json:"page"`
	Limit int        `json:"limit"`
	Total int        `json:"total"`
}

// GetAuditLogs godoc
//
//	@Summary		List audit log entries
//	@Description	Returns a paginated list of all audit log events. Filterable by username and event_type. Requires authentication.
//	@Tags			audit
//	@Produce		json
//	@Param			page		query	int		false	"Page number (default: 1)"
//	@Param			limit		query	int		false	"Results per page (default: 20, max: 100)"
//	@Param			username	query	string	false	"Filter by username"
//	@Param			event_type	query	string	false	"Filter by event type (e.g. LOGIN_FAILED, API_ACCESS)"
//	@Security		BearerAuth
//	@Success		200	{object}	PaginatedAuditLogs
//	@Failure		401	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/audit-logs [get]
func GetAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	username := c.Query("username")
	eventType := c.Query("event_type")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	args := []interface{}{}
	argIdx := 1
	where := "WHERE 1=1"

	if username != "" {
		where += fmt.Sprintf(" AND username = $%d", argIdx)
		args = append(args, username)
		argIdx++
	}
	if eventType != "" {
		where += fmt.Sprintf(" AND event_type = $%d", argIdx)
		args = append(args, eventType)
		argIdx++
	}

	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	if err := db.DB.QueryRow("SELECT COUNT(*) FROM audit_logs "+where, countArgs...).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	query := fmt.Sprintf(
		`SELECT id, event_type, COALESCE(username, ''), COALESCE(ip_address, ''), COALESCE(request_path, ''),
										 		COALESCE(request_method, ''), COALESCE(status_code, 0), timestamp::text, COALESCE(details, '')
										 		FROM audit_logs %s
										 		ORDER BY timestamp DESC
										 		LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	logs := []AuditLog{}
	for rows.Next() {
		var entry AuditLog
		if err := rows.Scan(
			&entry.ID, &entry.EventType, &entry.Username,
			&entry.IPAddress, &entry.RequestPath, &entry.RequestMethod,
			&entry.StatusCode, &entry.Timestamp, &entry.Details,
		); err != nil {
			continue
		}
		logs = append(logs, entry)
	}

	c.JSON(http.StatusOK, PaginatedAuditLogs{
		Data:  logs,
		Page:  page,
		Limit: limit,
		Total: total,
	})
}
