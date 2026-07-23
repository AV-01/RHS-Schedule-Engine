package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AV-01/RHS-Schedule-Engine/internal/db"
	"github.com/gin-gonic/gin"
)

type Student struct {
	ID            string `json: "id"`
	FirstName     string `json: "first_name"`
	MiddleInitial string `json: "middle_initial,omitempty"`
	LastName      string `json: "last_name"`
	Username      string `json: "username"`
}

type ScheduleEntry struct {
	SchoolYear  string `json: "school_year"`
	Grade       int    `json: "grade"`
	Period      int    `json: "period"`
	ClassName   string `json: "class_name"`
	TeacherName string `json: "teacher_name"`
	RoomNum     string `json: "room_num"`
}

type PaginatedStudents struct {
	Data  []Student `json: "data"`
	Page  int       `json: "page"`
	Limit int       `json: "limit"`
	Total int       `json: "total"`
}

// GetStudents godoc
//
//	@Summary		List all students
//	@Description	Returns a paginated list of students. Search by name or filter by grade.
//	@Tags			students
//	@Produce		json
//	@Param			page	query	int		false	"Page number (default: 1)"
//	@Param			limit	query	int		false	"Results per page (default: 20, max: 100)"
//	@Param			name	query	string	false	"Search by first or last name"
//	@Param			grade	query	int		false	"Filter by grade (9, 10, 11, 12)"
//	@Security		BearerAuth
//	@Success		200	{object}	PaginatedStudents
//	@Failure		401	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Router			/api/v1/students [get]
func GetStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	name := c.Query("name")
	grade := c.Query("grade")

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	if c.GetBool("is_demo") {
		c.JSON(http.StatusOK, GetDemoStudents(name, grade, page, limit))
		return
	}

	offset := (page - 1) * limit

	args := []interface{}{}
	argIdx := 1
	where := "WHERE 1=1"

	if name != "" {
		where += fmt.Sprintf(
			" AND (s.first_name ILIKE $%d OR s.last_name ILIKE $%d)",
			argIdx, argIdx,
		)

		args = append(args, "%"+name+"%")
		argIdx++
	}

	if grade != "" {
		where += fmt.Sprintf(" AND sc.grade = $%d", argIdx)
		args = append(args, grade)
		argIdx++
	}

	fromClause := "FROM students s LEFT JOIN schedules sc ON s.id = sc.student_uuid " + where

	selectQuery := fmt.Sprintf(`SELECT DISTINCT s.id::text, s.first_name, COALESCE(s.middle_initial, ''), s.last_name, s.username %s ORDER BY s.last_name, s.first_name LIMIT $%d OFFSET $%d`, fromClause, argIdx, argIdx+1)
	args = append(args, limit, offset)

	rows, err := db.DB.Query(selectQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.FirstName, &s.MiddleInitial, &s.LastName, &s.Username); err != nil {
			continue
		}
		students = append(students, s)
	}

	c.JSON(http.StatusOK, PaginatedStudents{
		Data:  students,
		Page:  page,
		Limit: limit,
		Total: len(students),
	})
}

// GetStudent godoc
//
//	@Summary		Get a student by UUID
//	@Description	Returns a single student's details by their internal UUID. Does not expose student ID.
//	@Tags			students
//	@Produce		json
//	@Param			id	path	string	true	"Student UUID"
//	@Security		BearerAuth
//	@Success		200	{object}	Student
//	@Failure		404	{object}	map[string]string
//	@Router			/api/v1/students/{id} [get]
func GetStudent(c *gin.Context) {
	id := c.Param("id")

	if c.GetBool("is_demo") {
		st, ok := GetDemoStudent(id)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
			return
		}
		c.JSON(http.StatusOK, st)
		return
	}

	var s Student
	err := db.DB.QueryRow(
		`SELECT id::text, first_name, COALESCE(middle_initial, ''), last_name, username
		FROM students WHERE id = $1`, id,
	).Scan(&s.ID, &s.FirstName, &s.MiddleInitial, &s.LastName, &s.Username)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	c.JSON(http.StatusOK, s)
}

// GetStudentSchedules godoc
//
//	@Summary		Get a student's schedules
//	@Description	Returns all schedule entries for a student across all years. Optionally filter by year.
//	@Tags			students
//	@Produce		json
//	@Param			id		path	string	true	"Student UUID"
//	@Param			year	query	string	false	"Filter by school year (e.g. 2023-2024)"
//	@Security		BearerAuth
//	@Success		200	{array}		ScheduleEntry
//	@Failure		404	{object}	map[string]string
//	@Router			/api/v1/students/{id}/schedules [get]
func GetStudentSchedules(c *gin.Context) {
	id := c.Param("id")
	year := c.Query("year")

	if c.GetBool("is_demo") {
		scheds := GetDemoStudentSchedules(id, year)
		if len(scheds) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "no schedules found for this studnet"})
			return
		}
		c.JSON(http.StatusOK, scheds)
		return
	}
	query := `
		SELECT sy.name, sc.grade, sc.period, sc.class_name, sc.teacher_name, sc.room_num
		FROM schedules sc
		JOIN school_years sy on sc.school_year_id = sy.id
		WHERE sc.student_uuid = $1
	`
	args := []interface{}{id}

	if year != "" {
		query += " AND sy.name = $2"
		args = append(args, year)
	}
	query += " ORDER BY sy.name, sc.period"

	rows, err := db.DB.Query(query, args...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	defer rows.Close()

	schedules := []ScheduleEntry{}
	for rows.Next() {
		var s ScheduleEntry
		if err := rows.Scan(&s.SchoolYear, &s.Grade, &s.Period, &s.ClassName, &s.TeacherName, &s.RoomNum); err != nil {
			continue
		}
		schedules = append(schedules, s)
	}

	if len(schedules) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no schedules found for this student"})
		return
	}

	c.JSON(http.StatusOK, schedules)
}
