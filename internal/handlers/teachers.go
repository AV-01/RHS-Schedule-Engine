package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
	"github.com/gin-gonic/gin"
)

type Teacher struct {
	Name string `json:"teacher_name"`
}

type TeacherScheduleEntry struct {
	SchoolYear string `json:"school_year"`
	Period     int    `json:"period"`
	ClassName  string `json:"class_name"`
	RoomNum    string `json:"room_num"`
}

type PaginatedTeachers struct {
	Data  []Teacher `json:"data"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
	Total int       `json:"total"`
}

// GetTeachers godoc
//
//	@Summary		List all teachers
//	@Description	Returns a paginated list of all unique teacher names across all school years.
//	@Tags			teachers
//	@Produce		json
//	@Param			page	query	int		false	"Page number (default: 1)"
//	@Param			limit	query	int		false	"Results per page (default: 20, limit: 100)"
//	@Param			name	query	string	false	"Search teacher name"
//	@Security		BearerAuth
//	@Success		200	{object}	PaginatedTeachers
//	@Failure		401	{object}	map[string]string
//	@Router			/api/v1/teachers [get]
func GetTeachers(c *gin.Context) {
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
	where := "WHERE teacher_name != ''"

	if name != "" {
		where += fmt.Sprintf(" AND teacher_name ILIKE $%d", argIdx)
		args = append(args, "%"+name+"%")
		argIdx++
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT teacher_name
		FROM schedules %s
		ORDER BY teacher_name
		LIMIT $%d OFFSET $%d
	`, where, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	defer rows.Close()

	teachers := []Teacher{}
	for rows.Next() {
		var t Teacher
		if err := rows.Scan(&t.Name); err != nil {
			continue
		}
		teachers = append(teachers, t)
	}

	c.JSON(http.StatusOK, PaginatedTeachers{
		Data:  teachers,
		Page:  page,
		Limit: limit,
		Total: len(teachers),
	})
}

// GetTeacherSchedule godoc
//
//	@Summary		Get a teacher's schedule
//	@Description	Returns all classes a teacher teaches, grouped by school year and period. Use the exact teacher_name value from GET /api/v1/teachers.
//	@Tags			teachers
//	@Produce		json
//	@Param			name	path	string	true	"Teacher name (e.g. Stark, L)"
//	@Security		BearerAuth
//	@Success		200	{array}		TeacherScheduleEntry
//	@Failure		404	{object}	map[string]string
//	@Router			/api/v1/teachers/{name}/schedules [get]
func GetTeacherSchedule(c *gin.Context) {
	name := c.Param("name")

	rows, err := db.DB.Query(
		`SELECT sy.name, sc.period, sc.class_name, sc.room_num
		FROM schedules sc
		JOIN school_years sy on sc.school_year_id = sy.id
		WHERE sc.teacher_name = $1
		ORDER BY sy.name, sc.period`,
		name,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	schedule := []TeacherScheduleEntry{}

	for rows.Next() {
		var entry TeacherScheduleEntry
		if err := rows.Scan(&entry.SchoolYear, &entry.Period, &entry.ClassName, &entry.RoomNum); err != nil {
			continue
		}
		schedule = append(schedule, entry)
	}

	if len(schedule) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no sched found, use GET /api/v1/teachers for valid names"})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
