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
		"templates/student/student.html",
		"templates/student/personal_account.html",
		"templates/student/my_group.html",
		"templates/student/schedule.html",
		"templates/student/discipline_progress.html",
		"templates/partials/grades-table.html",
		"templates/partials/error.html",
	))
}

// Основные страницы
// Главная страница - пока свободна
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "index", nil)
}

// Страница авторизации
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login", nil)
}

// Страница студента
func StudentHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "student", nil)
}

// Для HTMX
// Вход в приложение
func TryLogin(w http.ResponseWriter, r *http.Request) {
	// разобраться с возвратом кодов ошибок
	login := r.FormValue("login")
	pswd := r.FormValue("password")

	var user models.User
	result := config.DB.Where("login = ?", login).First(&user)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", "Неверный email или пароль")
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pswd))
	if err != nil {
		templates.ExecuteTemplate(w, "error", "Неверный email или пароль")
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
			templates.ExecuteTemplate(w, "error", "Проблемы на сервере, вернитесь позже")
			return
		}
	}
	cookie := &http.Cookie{
		Name:     "id_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,                   // Поменять на true
		SameSite: http.SameSiteStrictMode, // Защита от CSRF
	}
	http.SetCookie(w, cookie)
	// добавить редирект на остальные роли
	switch user.Role {
	case models.RoleStudent:
		w.Header().Set("HX-Redirect", "/student")
	default:
		w.Header().Set("HX-Redirect", "/")
	}
}

// Личный кабинет
func PersonalAccountHandler(w http.ResponseWriter, r *http.Request) {
	// Обо мне
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", "Проблемы на сервере, вернитесь позже")
	}
	var student models.Student
	result = config.DB.Where("id_user = ?", session.UserID).First(&student)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", "Проблемы на сервере, вернитесь позже")
	}
	var group models.Group
	result = config.DB.Where("id = ?", student.GroupID).First(&group)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", "Проблемы на сервере, вернитесь позже")
	}

	// Моя группа
	var students []models.Student
	result = config.DB.Where("id_group = ? AND id != ?", student.GroupID, student.ID).Find(&students)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", "Проблемы на сервере, вернитесь позже")
	}

	data := struct {
		LastName   string
		FirstName  string
		Patronymic string
		BirthDate  time.Time
		GroupSign  string
		Students   []models.Student
	}{
		student.LastName,
		student.FirstName,
		student.Patronymic,
		student.BirthDate,
		group.GroupSign,
		students,
	}

	templates.ExecuteTemplate(w, "personal_account", data)
}

// Расписание
func ScheduleHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "error", "Пока нет")
}

// Успеваемость
func DisciplineProgressHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "error", "Пока нет")
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
