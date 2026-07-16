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

	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	db.DB.QueryRow("SELECT COUNT(DISTINCT class_name) FROM schedules"+where, countArgs...).Scan(&total)

	query := fmt.Sprintf(
		`SELECT class_name, COUNT(*) as times_offered
		FROM schedules %s
		GROUP BY class_name
		ORDER BY class_name
		LIMIT $%d OFFSET $%d`,
		where, argIdx, argIdx+1,
	)

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
		Total: total,
	})
}
