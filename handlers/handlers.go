package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"html/template"
	"net/http"
	"opd_project/models"
	"time"
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
	// Создаём карту функций для шаблонов
	funcMap := template.FuncMap{
		// Арифметические функции
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},

		"div": func(a, b int) float64 {
			return float64(a) / float64(b)
		},

		// Функции для форматирования
		"formatDate": func(t time.Time) string {
			return t.Format("02.01.2006")
		},
		"formatDateTime": func(t time.Time) string {
			return t.Format("02.01.2006 15:04")
		},

		// Функции для работы с Action
		"getGradeValue": func(action models.Action) int {
			return action.Grade
		},
		"getAttendanceValue": func(action models.Action) string {
			attendanceMap := map[int]string{
				0: "Я",
				1: "Н",
				2: "Б",
				3: "ДО",
			}
			if val, ok := attendanceMap[int(action.Attendance)]; ok {
				return val
			}
			return ""
		},
		"getGradeDisplay": func(action models.Action) string {
			if action.Grade != 0 {
				return fmt.Sprintf("%d", action.Grade)
			}
			return "—"
		},
		"getAttendanceDisplay": func(action models.Action) string {
			attendanceMap := map[int]string{
				0: "Я",
				1: "Н",
				2: "Б",
				3: "ДО",
			}
			if val, ok := attendanceMap[int(action.Attendance)]; ok {
				return val
			}
			return "—"
		},
		"formatCell": func(action models.Action) string {
			if action.Grade != 0 {
				return fmt.Sprintf("%d", action.Grade)
			}
			attendanceMap := map[int]string{
				0: "Я",
				1: "Н",
				2: "Б",
				3: "ДО",
			}
			if val, ok := attendanceMap[int(action.Attendance)]; ok {
				return val
			}
			return "—"
		},

		// Функция для получения значения по индексу
		"getAction": func(actions [][]models.Action, i, j int) models.Action {
			if i < len(actions) && j < len(actions[i]) {
				return actions[i][j]
			}
			return models.Action{}
		},
	}

	// Создаём шаблон с функциями и парсим файлы
	templates = template.Must(
		template.New("").
			Funcs(funcMap).
			ParseFiles(
				"templates/index.html",
				"templates/login.html",
				"templates/student/student.html",
				"templates/student/personal_account.html",
				"templates/student/schedule.html",
				"templates/student/schedule_part.html",
				"templates/student/discipline_progress.html",
				"templates/student/dashboard.html",
				"templates/teacher/teacher.html",
				"templates/teacher/dashboard.html",
				"templates/teacher/personal_account.html",
				"templates/teacher/schedule.html",
				"templates/teacher/schedule_part.html",
				"templates/teacher/disciplines.html",
				"templates/teacher/disciplines_part_group.html",
				"templates/teacher/disciplines_part_table.html",
				"templates/tutor/tutor.html",
				"templates/error.html",
			),
	)
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
