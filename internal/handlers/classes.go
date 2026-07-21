package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
	"github.com/gin-gonic/gin"
)

type ClassEntry struct {
	ClassName string `json:"class_name"`
	Count     int    `json:"count"`
}
type PaginatedClasses struct {
	Data  []ClassEntry `json:"data"`
	Page  int          `json:"page"`
	Limit int          `json:"limit"`
	Total int          `json:"total"`
}

// GetClasses godoc
//
//	@Summary		List all unique classes
//	@Description	Returns a paginated list of all unique class names taught across all school years, with how many times each was offered.
//	@Tags			classes
//	@Produce		json
//	@Param			page	query	int		false	"Page number (default: 1)"
//	@Param			limit	query	int		false	"Results per page (default: 20, max: 100)"
//	@Param			name	query	string	false	"Search class name"
//	@Security		BearerAuth
//	@Success		200	{object}	PaginatedClasses
//	@Failure		401	{object}	map[string]string
//	@Router			/api/v1/classes [get]
func GetClasses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	name := c.Query("name")

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

	if name != "" {
		where += fmt.Sprintf(" AND class_name ILIKE $%d", argIdx)
		args = append(args, "%"+name+"%")
		argIdx++
	}

	query := fmt.Sprintf(
		`SELECT class_name, COUNT(*) as times_offered
		FROM schedules %s
		GROUP BY class_name
		ORDER BY class_name
		LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	classes := []ClassEntry{}
	for rows.Next() {
		var cls ClassEntry
		if err := rows.Scan(&cls.ClassName, &cls.Count); err != nil {
			continue
		}
		classes = append(classes, cls)
	}

	c.JSON(http.StatusOK, PaginatedClasses{
		Data:  classes,
		Page:  page,
		Limit: limit,
		Total: len(classes),
	})
}
