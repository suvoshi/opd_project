package models

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	Name   string
	Class  string
	Grades []Grade
}

type Grade struct {
	gorm.Model
	StudentID uint
	Value     int
	Subject   string
	Student   Student
}
