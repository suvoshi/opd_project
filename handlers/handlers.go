package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"time"
	// "strconv"
)

// Глобальная переменная для шаблонов (чтобы не компилировать их каждый раз)
var templates *template.Template

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// InitTemplates загружает все HTML шаблоны при старте
func InitTemplates() {
	templates = template.Must(template.ParseFiles(
		"templates/index.html",
		"templates/login.html",
		"templates/personal_account.html",
		"templates/my_group.html",
		"templates/schedule.html",
		"templates/discipline_progress.html",
		"templates/partials/grades-table.html",
	))
}

// Главная страница - Дашборд
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// cookie, err := r.Cookie("session_id")
	// if err != nil {
	// 	http.Redirect(w, r, "/login", http.StatusSeeOther)
	// 	return
	// }
	templates.ExecuteTemplate(w, "index", nil)
}

// Авторизация
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login", nil)
}

// Личный кабинет
func PersonalAccountHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "personal_account", nil)
}

// Моя группа
func MyGroupHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "my_group", nil)
}

// Расписание
func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "schedule", nil)
}

// Успеваемость
func DisciplineProgressHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "discipline_progress", nil)
}

// Получение таблицы оценок (для HTMX)
func GetLoginSession(w http.ResponseWriter, r *http.Request) {
	// TODO: сделать сессии через запрос к UserID (сессии для разных пользователей)
	sessionID := generateSessionID()
	session := models.Session{
		SessionID: sessionID,
		UserID:    1,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	config.DB.Create(&session)
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                   // Поменять на true
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
	}
	http.SetCookie(w, cookie)
}

func GetGradesHandler(w http.ResponseWriter, r *http.Request) {
	// var grades []models.Grade
	// // Preload подгружает связанные данные (имя ученика)
	// config.DB.Preload("Student").Find(&grades)

	// // Рендерим только фрагмент таблицы
	// templates.ExecuteTemplate(w, "grades-table.html", grades)
}

// Добавление оценки (для HTMX)
func AddGradeHandler(w http.ResponseWriter, r *http.Request) {
	// // Принимаем данные только через POST
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// // Парсим форму (аналог c.PostForm в Gin)
	// r.ParseForm()

	// studentID, _ := strconv.ParseUint(r.FormValue("student_id"), 10, 32)
	// value, _ := strconv.Atoi(r.FormValue("value"))
	// subject := r.FormValue("subject")

	// grade := models.Grade{
	// 	StudentID: uint(studentID),
	// 	Value:     value,
	// 	Subject:   subject,
	// }
	// config.DB.Create(&grade)

	// // Возвращаем обновленную таблицу (HTMX заменит её на странице)
	// GetGradesHandler(w, r)
}
