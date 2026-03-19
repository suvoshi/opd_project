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
}

type Student struct {
	ID         int    `gorm:"primaryKey"`
	UserID     int    `gorm:"column:id_user;not null;unique"`
	GroupID    int    `gorm:"column:id_group;not null"`
	LastName   string `gorm:"not null"`
	FirstName  string `gorm:"not null"`
	Patronymic string
	BirthDate  time.Time `gorm:"type:date;not null"`
}

type Lesson struct {
	ID           int `gorm:"primaryKey"`
	GroupID      int `gorm:"column:id_group;not null"`
	DisciplineID int `gorm:"column:id_discipline;not null"`
	Description  string
	DateBegin    time.Time `gorm:"not null"`
	DateEnd      time.Time `gorm:"not null"`
}

type Group struct {
	ID        int       `gorm:"primaryKey"`
	GroupSign string    `gorm:"unique"`
	DateBegin time.Time `gorm:"not null"`
	DateEnd   time.Time `gorm:"not null"`
}

type Discipline struct {
	ID          int    `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Hours       int
	TypePerform DisciplineTypePerform `gorm:"not null"`
}

type GroupDiscipline struct {
	ID           int       `gorm:"primaryKey"`
	GroupID      int       `gorm:"column:id_group;not null"`
	DisciplineID int       `gorm:"column:id_discipline;not null"`
	DateBegin    time.Time `gorm:"type:date;not null"`
	DateEnd      time.Time `gorm:"type:date;not null"`
}

type Action struct {
	ID         int `gorm:"primaryKey"`
	LessonID   int `gorm:"column:id_lesson;not null"`
	StudentID  int `gorm:"column:id_student;not null"`
	Grade      int `gorm:"check:grade BETWEEN 2 AND 5"`
	Attendance AttendanceType
}

type StudentEndDiscipline struct {
	ID                int       `gorm:"primaryKey"`
	GroupDisciplineID int       `gorm:"column:id_group_discipline;not null"`
	StudentID         int       `gorm:"column:id_student;not null"`
	Grade             int       `gorm:"not null"`
	Date              time.Time `gorm:"type:date;not null"`
}

type Session struct {
	SessionID string    `gorm:"column:id_session;not null;unique"`
	UserID    int       `gorm:"column:id_user;not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
}
