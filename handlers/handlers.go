package handlers

import (
	"html/template"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"strconv"
)

// Глобальная переменная для шаблонов (чтобы не компилировать их каждый раз)
var templates *template.Template

// InitTemplates загружает все HTML шаблоны при старте
func InitTemplates() {
	templates = template.Must(template.ParseGlob("templates/**/*.html"))
}

// Главная страница
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Рендерим шаблон с именем "base.html"
	templates.ExecuteTemplate(w, "base.html", nil)
}

// Получение таблицы оценок (для HTMX)
func GetGradesHandler(w http.ResponseWriter, r *http.Request) {
	var grades []models.Grade
	// Preload подгружает связанные данные (имя ученика)
	config.DB.Preload("Student").Find(&grades)

	// Рендерим только фрагмент таблицы
	templates.ExecuteTemplate(w, "grades-table.html", grades)
}

// Добавление оценки (для HTMX)
func AddGradeHandler(w http.ResponseWriter, r *http.Request) {
	// Принимаем данные только через POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Парсим форму (аналог c.PostForm в Gin)
	r.ParseForm()

	studentID, _ := strconv.ParseUint(r.FormValue("student_id"), 10, 32)
	value, _ := strconv.Atoi(r.FormValue("value"))
	subject := r.FormValue("subject")

	grade := models.Grade{
		StudentID: uint(studentID),
		Value:     value,
		Subject:   subject,
	}
	config.DB.Create(&grade)

	// Возвращаем обновленную таблицу (HTMX заменит её на странице)
	GetGradesHandler(w, r)
}
