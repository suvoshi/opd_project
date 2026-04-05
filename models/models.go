package models

import (
	"time"
)

// ENUM типы
type UserRole int
type DisciplineTypePerform int
type AttendanceType int

const (
	// Роли для пользователей
	RoleStudent UserRole = 0
	RoleTeacher UserRole = 1
	RoleTutor   UserRole = 2
	RoleAdmin   UserRole = 3

	// Варианты проведения занятий
	Lection  DisciplineTypePerform = 0
	Practice DisciplineTypePerform = 1

	// Посещаемость, явка, неявка, отсутствие по болезни и ДО соотсветственно
	Present AttendanceType = 0
	Absent  AttendanceType = 1
	Sick    AttendanceType = 2
	DO      AttendanceType = 3
)

type User struct {
	ID       int      `gorm:"primaryKey"`
	Login    string   `gorm:"unique;not null"`
	Password string   `gorm:"not null"`
	Role     UserRole `gorm:"not null"`

	Student *Student `gorm:"foreignKey:UserID"`
	Teacher *Teacher `gorm:"foreignKey:UserID"`
	Tutor   *Tutor   `gorm:"foreignKey:UserID"`
}

type Student struct {
	ID         int    `gorm:"primaryKey"`
	UserID     int    `gorm:"column:id_user;not null;unique"`
	GroupID    int    `gorm:"column:id_group;not null"`
	LastName   string `gorm:"not null"`
	FirstName  string `gorm:"not null"`
	Patronymic string
	BirthDate  time.Time `gorm:"type:date;not null"`

	User    User     `gorm:"foreignKey:UserID"`
	Group   Group    `gorm:"foreignKey:GroupID"`
	Actions []Action `gorm:"foreignKey:StudentID"`
}

type Lesson struct {
	ID           int `gorm:"primaryKey"`
	GroupID      int `gorm:"column:id_group;not null"`
	DisciplineID int `gorm:"column:id_discipline;not null"`
	Description  string
	DateBegin    time.Time `gorm:"not null"`
	DateEnd      time.Time `gorm:"not null"`
	TeacherID    int       `gorm:"column:id_teacher"`

	Group      Group      `gorm:"foreignKey:GroupID"`
	Discipline Discipline `gorm:"foreignKey:DisciplineID"`
	Actions    []Action   `gorm:"foreignKey:LessonID"`
	Teacher    Teacher    `gorm:"foreignKey:TeacherID"`
}

type Group struct {
	ID        int       `gorm:"primaryKey"`
	GroupSign string    `gorm:"unique"`
	DateBegin time.Time `gorm:"not null"`
	DateEnd   time.Time `gorm:"not null"`

	Students         []Student         `gorm:"foreignKey:GroupID"`
	Lessons          []Lesson          `gorm:"foreignKey:GroupID"`
	GroupDisciplines []GroupDiscipline `gorm:"foreignKey:GroupID"`
}

type Discipline struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Hours       int
	TypePerform DisciplineTypePerform `gorm:"not null"`

	Lessons          []Lesson          `gorm:"foreignKey:DisciplineID"`
	GroupDisciplines []GroupDiscipline `gorm:"foreignKey:DisciplineID"`
}

type GroupDiscipline struct {
	ID           int       `gorm:"primaryKey"`
	GroupID      int       `gorm:"column:id_group;not null"`
	DisciplineID int       `gorm:"column:id_discipline;not null"`
	DateBegin    time.Time `gorm:"type:date;not null"`
	DateEnd      time.Time `gorm:"type:date;not null"`
	TeacherID    int       `gorm:"column:id_teacher"`

	Group      Group      `gorm:"foreignKey:GroupID"`
	Discipline Discipline `gorm:"foreignKey:DisciplineID"`
	Teacher    Teacher    `gorm:"foreignKey:TeacherID"`
}

type Action struct {
	ID         int `gorm:"primaryKey"`
	LessonID   int `gorm:"column:id_lesson;not null"`
	StudentID  int `gorm:"column:id_student;not null"`
	Grade      int `gorm:"check:grade BETWEEN 2 AND 5"`
	Attendance AttendanceType
	// тот, кто поставил grade или attendance (не только учитель)
	UserID int `gorm:"column:id_user"`

	Lesson  Lesson  `gorm:"foreignKey:LessonID"`
	Student Student `gorm:"foreignKey:StudentID"`
	User    User    `gorm:"foreignKey:UserID"`
}

type StudentEndDiscipline struct {
	ID                int       `gorm:"primaryKey"`
	GroupDisciplineID int       `gorm:"column:id_group_discipline;not null"`
	StudentID         int       `gorm:"column:id_student;not null"`
	Grade             int       `gorm:"not null"`
	Date              time.Time `gorm:"type:date;not null"`
	// тот, кто аттестовал студента по дисциплине (не только учитель)
	UserID int `gorm:"column:id_user"`

	GroupDiscipline GroupDiscipline `gorm:"foreignKey:GroupDisciplineID"`
	Student         Student         `gorm:"foreignKey:StudentID"`
	User            User            `gorm:"foreignKey:UserID"`
}

type Session struct {
	SessionID string    `gorm:"column:id_session;not null;unique"`
	UserID    int       `gorm:"column:id_user;not null"`
	ExpiresAt time.Time `gorm:"not null"`

	User User `gorm:"foreignKey:UserID"`
}

type Teacher struct {
	ID         int    `gorm:"primaryKey"`
	UserID     int    `gorm:"column:id_user;not null;unique"`
	LastName   string `gorm:"not null"`
	FirstName  string `gorm:"not null"`
	Patronymic string
	BirthDate  time.Time `gorm:"type:date;not null"`

	User             User              `gorm:"foreignKey:UserID"`
	GroupDisciplines []GroupDiscipline `gorm:"foreignKey:TeacherID"`
	Lessons          []Lesson          `gorm:"foreignKey:TeacherID"`
}

type Tutor struct {
	ID         int    `gorm:"primaryKey"`
	UserID     int    `gorm:"column:id_user;not null;unique"`
	LastName   string `gorm:"not null"`
	FirstName  string `gorm:"not null"`
	Patronymic string
	BirthDate  time.Time `gorm:"type:date;not null"`

	User User `gorm:"foreignKey:UserID"`
}
