package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"sort"
	"time"
	//"strconv"
)

// HTMX хендлеры для студента

// Данные, используемые хендлерами для заполнения шаблонов
type StudentDashboardData struct {
	GradeData []struct {
		DisciplineName string
		Grade          int
		Date           time.Time
	}
	AnnouncementData []models.Announcement
}

type StudentPersonalAccountData struct {
	Student  models.Student
	Students []models.Student
}

type StudentSchedulePartData struct {
	WeekLessons [][]models.Lesson
}

type StudentDisciplineProgressData struct {
	DisciplineData []struct {
		Discipline models.Discipline
		Actions    []models.Action
		FinalGrade string
	}
}

// Дашборд
func StudentDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
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
		return
	}
	slog.Info("StudentDashboardHandler - Пытаемся найти студента", "id_user", session.UserID)
	var student models.Student
	result = config.DB.Where("id_user = ?", session.UserID).First(&student)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	now := time.Now()
	weekAgo := now.Add(-7 * 24 * time.Hour)

	data := StudentDashboardData{}

	result = config.DB.
		Table("actions").
		Select("disciplines.name as discipline_name, actions.grade, actions.created_at as date").
		Joins("JOIN lessons ON lessons.id = actions.id_lesson").
		Joins("JOIN disciplines ON disciplines.id = lessons.id_discipline").
		Where("actions.id_student = ?", student.ID).
		Where("actions.created_at BETWEEN ? AND ?", weekAgo, now).
		Order("actions.created_at DESC").
		Scan(&data.GradeData)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	result = config.DB.
		Where("(date BETWEEN ? AND ?) AND visibility <= 1", weekAgo, now).
		Order("date DESC").
		Find(&data.AnnouncementData)

	templates.ExecuteTemplate(w, "dashboard", data)
	slog.Info("StudentDashboardHandler - Успешно", "id_user", session.UserID)
}

// Личный кабинет
func StudentPersonalAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
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
		return
	}
	slog.Info("StudentPersonalAccountHandler - Пытаемся найти студента", "id_user", session.UserID)
	var student models.Student
	result = config.DB.
		Preload("Group").
		Where("id_user = ?", session.UserID).
		First(&student)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	// Группа (все, кроме студента, сделавшего запрос)
	var students []models.Student
	result = config.DB.Where("id_group = ? AND id != ?", student.GroupID, student.ID).Find(&students)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	data := StudentPersonalAccountData{Student: student, Students: students}

	templates.ExecuteTemplate(w, "personal_account", data)
	slog.Info("StudentPersonalAccountHandler - Успешно", "id_user", session.UserID)
}

// Расписание (выбор недели)
func StudentScheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	templates.ExecuteTemplate(w, "schedule", nil)
}

// Расписание (время проведения занятий выбранной недели)
func StudentSchedulePartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
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
		return
	}
	slog.Info("StudentSchedulePartHandler - Пытаемся найти студента", "id_user", session.UserID)
	var student models.Student
	result = config.DB.Where("id_user = ?", session.UserID).First(&student)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
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

	data := StudentSchedulePartData{WeekLessons: weekLessons}

	templates.ExecuteTemplate(w, "schedule_part", data)
	slog.Info("StudentSchedulePartHandler - Успешно", "id_user", session.UserID)
}

// Успеваемость (итоговые оценки по предметам)
func StudentDisciplineProgressHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
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
		return
	}
	slog.Info("StudentDisciplineProgressHandler - Пытаемся найти студента", "id_user", session.UserID)
	var student models.Student
	result = config.DB.Where("id_user = ?", session.UserID).First(&student)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	now := time.Now()

	var groupDisciplines []models.GroupDiscipline
	result = config.DB.
		Preload("Discipline").
		Where("(? BETWEEN date_begin AND date_end) AND id_group = ?", now, student.GroupID).
		Find(&groupDisciplines)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	sort.Slice(groupDisciplines, func(i, j int) bool {
		return groupDisciplines[i].Discipline.Name < groupDisciplines[j].Discipline.Name
	})

	data := StudentDisciplineProgressData{}

	for _, groupDisc := range groupDisciplines {
		var actions []models.Action
		result = config.DB.
			Joins("JOIN lessons ON lessons.id = actions.id_lesson").
			Where("actions.id_student = ?", student.ID).
			Where("lessons.id_discipline = ?", groupDisc.DisciplineID).
			Find(&actions)
		if result.Error != nil {
			templates.ExecuteTemplate(w, "error", errorServerSide)
			return
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

		data.DisciplineData = append(data.DisciplineData, struct {
			Discipline models.Discipline
			Actions    []models.Action
			FinalGrade string
		}{
			Discipline: groupDisc.Discipline,
			Actions:    actions,
			FinalGrade: finalGrade,
		})
	}
	templates.ExecuteTemplate(w, "discipline_progress", data)
	slog.Info("StudentDisciplineProgressHandler - Успешно", "id_user", session.UserID)
}
