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

	var total int
	countArgs := make([]interface{}, len(args))
	copy(countArgs, args)
	if err := db.DB.QueryRow("SELECT COUNT(DISTINCT s.id) "+fromClause, countArgs...).Scan(&total); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

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
		Total: total,
	})
}

func GetStudent(c *gin.Context) {
	id := c.Param("id")

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

func GetStudentSchedules(c *gin.Context) {
	id := c.Param("id")
	year := c.Query("year")

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
