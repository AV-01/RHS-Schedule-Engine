package handlers

import (
	"strconv"
	"strings"
)

type demoStudent struct {
	Student
	Grade int
}

var mockStudents = []demoStudent{
	{
		Student: Student{
			ID:            "demo-student-001-uuid",
			FirstName:     "Alex",
			MiddleInitial: "J",
			LastName:      "Morgan",
			Username:      "alex.morgan",
		},
		Grade: 9,
	},
	{
		Student: Student{
			ID:            "demo-student-002-uuid",
			FirstName:     "Jordan",
			MiddleInitial: "K",
			LastName:      "Smith",
			Username:      "jordan.smith",
		},
		Grade: 10,
	},
	{
		Student: Student{
			ID:            "demo-student-003-uuid",
			FirstName:     "Taylor",
			MiddleInitial: "R",
			LastName:      "Swift",
			Username:      "taylor.swift",
		},
		Grade: 11,
	},
	{
		Student: Student{
			ID:            "demo-student-004-uuid",
			FirstName:     "Sam",
			MiddleInitial: "L",
			LastName:      "Wilson",
			Username:      "sam.wilson",
		},
		Grade: 12,
	},
	{
		Student: Student{
			ID:            "demo-student-005-uuid",
			FirstName:     "Casey",
			MiddleInitial: "M",
			LastName:      "Neistat",
			Username:      "casey.neistat",
		},
		Grade: 11,
	},
}

var mockStudentSchedules = map[string][]ScheduleEntry{
	"demo-student-001-uuid": {
		{SchoolYear: "2023-2024", Grade: 9, Period: 1, ClassName: "AP Computer Science A", TeacherName: "Stark, T", RoomNum: "C101"},
		{SchoolYear: "2023-2024", Grade: 9, Period: 2, ClassName: "Honors Chemistry", TeacherName: "Curie, M", RoomNum: "S204"},
		{SchoolYear: "2023-2024", Grade: 9, Period: 3, ClassName: "English 9", TeacherName: "Shakespeare, W", RoomNum: "E105"},
		{SchoolYear: "2023-2024", Grade: 9, Period: 4, ClassName: "Algebra II", TeacherName: "Euler, L", RoomNum: "M302"},
	},
	"demo-student-002-uuid": {
		{SchoolYear: "2023-2024", Grade: 10, Period: 1, ClassName: "US History", TeacherName: "Washington, G", RoomNum: "H201"},
		{SchoolYear: "2023-2024", Grade: 10, Period: 2, ClassName: "AP Computer Science A", TeacherName: "Stark, T", RoomNum: "C101"},
		{SchoolYear: "2023-2024", Grade: 10, Period: 3, ClassName: "Physical Education", TeacherName: "Johnson, M", RoomNum: "GYM"},
		{SchoolYear: "2023-2024", Grade: 10, Period: 4, ClassName: "Algebra II", TeacherName: "Euler, L", RoomNum: "M302"},
	},
	"demo-student-003-uuid": {
		{SchoolYear: "2023-2024", Grade: 11, Period: 1, ClassName: "AP Physics C", TeacherName: "Einstein, A", RoomNum: "S301"},
		{SchoolYear: "2023-2024", Grade: 11, Period: 2, ClassName: "Honors English 11", TeacherName: "Shakespeare, W", RoomNum: "E105"},
		{SchoolYear: "2023-2024", Grade: 11, Period: 3, ClassName: "AP Chemistry", TeacherName: "Curie, M", RoomNum: "S204"},
		{SchoolYear: "2023-2024", Grade: 11, Period: 4, ClassName: "Physical Education", TeacherName: "Johnson, M", RoomNum: "GYM"},
	},
	"demo-student-004-uuid": {
		{SchoolYear: "2023-2024", Grade: 12, Period: 1, ClassName: "AP Computer Science A", TeacherName: "Stark, T", RoomNum: "C101"},
		{SchoolYear: "2023-2024", Grade: 12, Period: 2, ClassName: "AP US History", TeacherName: "Washington, G", RoomNum: "H201"},
		{SchoolYear: "2023-2024", Grade: 12, Period: 3, ClassName: "AP Physics C", TeacherName: "Einstein, A", RoomNum: "S301"},
		{SchoolYear: "2023-2024", Grade: 12, Period: 4, ClassName: "Intro to Robotics", TeacherName: "Stark, T", RoomNum: "C101"},
	},
	"demo-student-005-uuid": {
		{SchoolYear: "2023-2024", Grade: 11, Period: 1, ClassName: "Algebra II", TeacherName: "Euler, L", RoomNum: "M302"},
		{SchoolYear: "2023-2024", Grade: 11, Period: 2, ClassName: "Honors Chemistry", TeacherName: "Curie, M", RoomNum: "S204"},
		{SchoolYear: "2023-2024", Grade: 11, Period: 3, ClassName: "English 9", TeacherName: "Shakespeare, W", RoomNum: "E105"},
		{SchoolYear: "2023-2024", Grade: 11, Period: 4, ClassName: "Intro to Robotics", TeacherName: "Stark, T", RoomNum: "C101"},
	},
}

var mockClasses = []ClassEntry{
	{ClassName: "Algebra II", Count: 185},
	{ClassName: "AP Chemistry", Count: 88},
	{ClassName: "AP Computer Science A", Count: 142},
	{ClassName: "AP Physics C", Count: 95},
	{ClassName: "AP US History", Count: 110},
	{ClassName: "English 9", Count: 230},
	{ClassName: "Honors Chemistry", Count: 118},
	{ClassName: "Honors English 11", Count: 135},
	{ClassName: "Intro to Robotics", Count: 76},
	{ClassName: "Physical Education", Count: 310},
	{ClassName: "US History", Count: 190},
}

var mockTeachers = []Teacher{
	{Name: "Curie, M"},
	{Name: "Einstein, A"},
	{Name: "Euler, L"},
	{Name: "Johnson, M"},
	{Name: "Shakespeare, W"},
	{Name: "Stark, T"},
	{Name: "Washington, G"},
}

var mockTeacherSchedules = map[string][]TeacherScheduleEntry{
	"Stark, T": {
		{SchoolYear: "2023-2024", Period: 1, ClassName: "AP Computer Science A", RoomNum: "C101"},
		{SchoolYear: "2023-2024", Period: 2, ClassName: "AP Computer Science A", RoomNum: "C101"},
		{SchoolYear: "2023-2024", Period: 4, ClassName: "Intro to Robotics", RoomNum: "C101"},
	},
	"Curie, M": {
		{SchoolYear: "2023-2024", Period: 2, ClassName: "Honors Chemistry", RoomNum: "S204"},
		{SchoolYear: "2023-2024", Period: 3, ClassName: "AP Chemistry", RoomNum: "S204"},
	},
	"Einstein, A": {
		{SchoolYear: "2023-2024", Period: 1, ClassName: "AP Physics C", RoomNum: "S301"},
		{SchoolYear: "2023-2024", Period: 3, ClassName: "AP Physics C", RoomNum: "S301"},
	},
	"Euler, L": {
		{SchoolYear: "2023-2024", Period: 1, ClassName: "Algebra II", RoomNum: "M302"},
		{SchoolYear: "2023-2024", Period: 4, ClassName: "Algebra II", RoomNum: "M302"},
	},
	"Johnson, M": {
		{SchoolYear: "2023-2024", Period: 3, ClassName: "Physical Education", RoomNum: "GYM"},
		{SchoolYear: "2023-2024", Period: 4, ClassName: "Physical Education", RoomNum: "GYM"},
	},
	"Shakespeare, W": {
		{SchoolYear: "2023-2024", Period: 2, ClassName: "Honors English 11", RoomNum: "E105"},
		{SchoolYear: "2023-2024", Period: 3, ClassName: "English 9", RoomNum: "E105"},
	},
	"Washington, G": {
		{SchoolYear: "2023-2024", Period: 1, ClassName: "US History", RoomNum: "H201"},
		{SchoolYear: "2023-2024", Period: 2, ClassName: "AP US History", RoomNum: "H201"},
	},
}

// GetDemoStudents returns paginated mock students based on name and grade filters
func GetDemoStudents(name, grade string, page, limit int) PaginatedStudents {
	var filtered []Student
	targetGrade, _ := strconv.Atoi(grade)

	for _, s := range mockStudents {
		if name != "" {
			queryLower := strings.ToLower(name)
			firstLower := strings.ToLower(s.FirstName)
			lastLower := strings.ToLower(s.LastName)
			if !strings.Contains(firstLower, queryLower) && !strings.Contains(lastLower, queryLower) {
				continue
			}
		}
		if grade != "" && s.Grade != targetGrade {
			continue
		}
		filtered = append(filtered, s.Student)
	}

	total := len(filtered)
	offset := (page - 1) * limit
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}

	paged := filtered[offset:end]
	if paged == nil {
		paged = []Student{}
	}

	return PaginatedStudents{
		Data:  paged,
		Page:  page,
		Limit: limit,
		Total: total,
	}
}

// GetDemoStudent returns a mock student by UUID or default demo student if not found
func GetDemoStudent(id string) (Student, bool) {
	for _, s := range mockStudents {
		if s.ID == id {
			return s.Student, true
		}
	}
	// Fallback to first mock student for demo requests
	if len(mockStudents) > 0 {
		st := mockStudents[0].Student
		st.ID = id
		return st, true
	}
	return Student{}, false
}

// GetDemoStudentSchedules returns mock schedules for a student, optionally filtered by year
func GetDemoStudentSchedules(id, year string) []ScheduleEntry {
	scheds, ok := mockStudentSchedules[id]
	if !ok {
		// Default mock schedules for demo fallback
		scheds = mockStudentSchedules["demo-student-001-uuid"]
	}

	if year == "" {
		return scheds
	}

	var filtered []ScheduleEntry
	for _, sc := range scheds {
		if sc.SchoolYear == year {
			filtered = append(filtered, sc)
		}
	}
	return filtered
}

// GetDemoClasses returns paginated mock classes based on name query
func GetDemoClasses(name string, page, limit int) PaginatedClasses {
	var filtered []ClassEntry

	for _, cls := range mockClasses {
		if name != "" {
			if !strings.Contains(strings.ToLower(cls.ClassName), strings.ToLower(name)) {
				continue
			}
		}
		filtered = append(filtered, cls)
	}

	total := len(filtered)
	offset := (page - 1) * limit
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}

	paged := filtered[offset:end]
	if paged == nil {
		paged = []ClassEntry{}
	}

	return PaginatedClasses{
		Data:  paged,
		Page:  page,
		Limit: limit,
		Total: total,
	}
}

// GetDemoTeachers returns paginated mock teachers based on name query
func GetDemoTeachers(name string, page, limit int) PaginatedTeachers {
	var filtered []Teacher

	for _, t := range mockTeachers {
		if name != "" {
			if !strings.Contains(strings.ToLower(t.Name), strings.ToLower(name)) {
				continue
			}
		}
		filtered = append(filtered, t)
	}

	total := len(filtered)
	offset := (page - 1) * limit
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}

	paged := filtered[offset:end]
	if paged == nil {
		paged = []Teacher{}
	}

	return PaginatedTeachers{
		Data:  paged,
		Page:  page,
		Limit: limit,
		Total: total,
	}
}

// GetDemoTeacherSchedule returns mock teacher schedules by teacher name
func GetDemoTeacherSchedule(name string) []TeacherScheduleEntry {
	scheds, ok := mockTeacherSchedules[name]
	if !ok {
		return []TeacherScheduleEntry{}
	}
	return scheds
}
