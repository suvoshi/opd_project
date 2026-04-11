package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"sort"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	//"strconv"
)

// Сообщение об ошибках, потом подкорректировать
var errorServerSide = "Проблемы на сервере, вернитесь позже"
var incorrectEmailOrLogin = "Неправильный email или логин"

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
		"templates/student/schedule.html",
		"templates/student/schedule_part.html",
		"templates/student/discipline_progress.html",
		"templates/student/dashboard.html",
		"templates/teacher/teacher.html",
		"templates/tutor/tutor.html",
		"templates/error.html",
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

// Страница преподователя
func TeacherHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "teacher", nil)
}

// Страница куратора
func TutorHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "tutor", nil)
}

// Для HTMX
// Вход в приложение
func TryLogin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		// потом сделать обработку такого запроса
	}
	login := r.FormValue("login")
	pswd := r.FormValue("password")

	var user models.User
	result := config.DB.Where("login = ?", login).First(&user)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", incorrectEmailOrLogin)
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pswd))
	if err != nil {
		templates.ExecuteTemplate(w, "error", incorrectEmailOrLogin)
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
			templates.ExecuteTemplate(w, "error", errorServerSide)
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
		MaxAge:   86400,
	}
	http.SetCookie(w, cookie)

	switch user.Role {
	case models.RoleStudent:
		w.Header().Set("HX-Redirect", "/student")
	case models.RoleTeacher:
		w.Header().Set("HX-Redirect", "/teacher")
	case models.RoleTutor:
		w.Header().Set("HX-Redirect", "/tutor")
	case models.RoleAdmin:
		w.Header().Set("HX-Redirect", "/admin")
	}
	slog.Info("TryLogin - Успешный вход",
		"id_user", user.ID, "role", user.Role)
}

// Дашборд
func StudentDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		// потом сделать обработку такого запроса
	}
	templates.ExecuteTemplate(w, "dashboard", nil)
}

// Личный кабинет
func StudentPersonalAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		// потом сделать обработку такого запроса
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// Получение данных о студенте
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
	}
	slog.Info("StudentPersonalAccountHandler - Пытаемся найти студента", "id_user", session.UserID)
	var student models.Student
	result = config.DB.
		Preload("Group").
		Where("id_user = ?", session.UserID).
		First(&student)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
	}

	// Группа (все, кроме студента, сделавшего запрос)
	var students []models.Student
	result = config.DB.Where("id_group = ? AND id != ?", student.GroupID, student.ID).Find(&students)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
	}

	// Отдаем ответ
	data := struct {
		Student  models.Student
		Students []models.Student
	}{
		Student:  student,
		Students: students,
	}

	templates.ExecuteTemplate(w, "personal_account", data)
	slog.Info("StudentPersonalAccountHandler - Успешно", "id_user", session.UserID)
}

// Расписание
func StudentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		// потом сделать обработку такого запроса
	}
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "schedule", nil)
}

// Расписание - по дням недели
func StudentSchedulePartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		// потом сделать обработку такого запроса
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим студента
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
	}
	slog.Info("StudentSchedulePartHandler - Пытаемся найти студента", "id_user", session.UserID)
	var student models.Student
	result = config.DB.Where("id_user = ?", session.UserID).First(&student)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
	}

	r.ParseForm()
	// всегда понедельник и воскресенье соответственно
	start, _ := time.Parse("2006-01-02", r.FormValue("start"))
	//end, _ := time.Parse("2006-01-02", r.FormValue("end"))

	weekLessons := make([][]models.Lesson, 7)
	point1 := start
	point2 := start.Add(24 * time.Hour)

	for ind := 0; ind < 7; ind++ {
		var lessons []models.Lesson
		result = config.DB.
			Preload("Discipline").
			Preload("Actions", "id_student = ?", student.ID).
			Where("id_group = ? AND (? < date_begin AND date_end < ?)", student.GroupID, point1, point2).
			Order("date_begin").
			Find(&lessons)
		if result.Error != nil {
			templates.ExecuteTemplate(w, "error", errorServerSide)
			return
		}
		weekLessons[ind] = lessons
		point1, point2 = point2, point2.Add(24*time.Hour)
	}

	data := struct {
		WeekLessons [][]models.Lesson
	}{
		WeekLessons: weekLessons,
	}
	templates.ExecuteTemplate(w, "schedule_part", data)
	slog.Info("StudentSchedulePartHandler - Успешно", "id_user", session.UserID)
}

// Успеваемость
func StudentDisciplineProgressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		// потом сделать обработку такого запроса
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим студента
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
	}
	slog.Info("StudentDisciplineProgressHandler - Пытаемся найти студента", "id_user", session.UserID)
	var student models.Student
	result = config.DB.Where("id_user = ?", session.UserID).First(&student)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
	}

	now := time.Now()

	var groupDisciplines []models.GroupDiscipline
	result = config.DB.
		Preload("Discipline").
		Where("(? BETWEEN date_begin AND date_end) AND id_group = ?", now, student.GroupID).
		Find(&groupDisciplines)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
	}
	sort.Slice(groupDisciplines, func(i, j int) bool {
		return groupDisciplines[i].Discipline.Name < groupDisciplines[j].Discipline.Name
	})

	data := make([]struct {
		Discipline models.Discipline
		Actions    []models.Action
		FinalGrade string
	}, len(groupDisciplines))

	for ind, groupDisc := range groupDisciplines {
		var actions []models.Action
		result = config.DB.
			Joins("JOIN lessons ON lessons.id = actions.id_lesson").
			Where("actions.id_student = ?", student.ID).
			Where("lessons.id_discipline = ?", groupDisc.DisciplineID).
			Find(&actions)
		if result.Error != nil {
			templates.ExecuteTemplate(w, "error", errorServerSide)
		}
		count := 0.0
		sum := 0.0

		for _, action := range actions {
			count += 1
			sum += float64(action.Grade)
		}

		var finalGrade string
		if count == 0 {
			finalGrade = "-"
		} else {
			finalGrade = fmt.Sprintf("%.2f", sum/count)
		}

		data[ind] = struct {
			Discipline models.Discipline
			Actions    []models.Action
			FinalGrade string
		}{
			Discipline: groupDisc.Discipline,
			Actions:    actions,
			FinalGrade: finalGrade,
		}
	}
	templates.ExecuteTemplate(w, "discipline_progress", data)
	slog.Info("StudentDisciplineProgressHandler - Успешно", "id_user", session.UserID)
}
