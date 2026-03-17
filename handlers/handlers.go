package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"html/template"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
	_, err := r.Cookie("session_id")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
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

// Для HTMX
// Вход в приложение
func GetLoginSession(w http.ResponseWriter, r *http.Request) {
	// разобраться с возвратом кодов ошибок
	login := r.FormValue("login")
	pswd := r.FormValue("password")

	var user models.User
	result := config.DB.Where("login = ?", login).First(&user)
	if result.Error != nil {
		// здесь и далее
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<div class="error">Неверный email или пароль</div>`))
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pswd))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<div class="error">Неверный email или пароль</div>`))
		return
	}
	sessionID := generateSessionID()
	session := models.Session{
		SessionID: sessionID,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	result = config.DB.Create(&session)
	for result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			sessionID = generateSessionID()
			session.SessionID = sessionID
			result = config.DB.Create(&session)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`<div class="error">Ошибка сервера</div>`))
			return
		}
	}
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                   // Поменять на true
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
	}
	http.SetCookie(w, cookie)
	w.Header().Set("HX-Redirect", "/")
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
