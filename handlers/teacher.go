package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"strconv"
	"time"

	"gorm.io/gorm"
	//"strconv"
)

// HTMX хендлеры для учителя

// Данные, используемые хендлерами для заполнения шаблонов
type TeacherDashboardData struct {
	AnnouncementData []models.Announcement
}

type TeacherPersonalAccountData struct {
	Teacher models.Teacher
}

type TeacherSchedulePartData struct {
	WeekLessons [][]models.Lesson
}

type TeacherDisciplinesData struct {
	GroupDisciplines []models.GroupDiscipline
}

type TeacherDisciplinesPartGroupData struct {
	GroupDisciplines []models.GroupDiscipline
}

type TeacherDisciplinesPartTableData struct {
	Group      models.Group
	Discipline models.Discipline
	Students   []models.Student
	Lessons    []models.Lesson
	Actions    [][]models.Action
}

// Дашборд
func TeacherDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим учителя
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TeacherDashboardHandler - Пытаемся найти учителя", "id_user", session.UserID)
	var teacher models.Teacher
	result = config.DB.Where("id_user = ?", session.UserID).First(&teacher)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	now := time.Now()
	weekAgo := now.Add(-7 * 24 * time.Hour)

	data := TeacherDashboardData{}

	result = config.DB.
		Where("(date BETWEEN ? AND ?) AND visibility <= 1", weekAgo, now).
		Order("date DESC").
		Find(&data.AnnouncementData)

	templates.ExecuteTemplate(w, "teacher_dashboard", data)
	slog.Info("TeacherDashboardHandler - Успешно", "id_user", session.UserID)
}

// Личный кабинет
func TeacherPersonalAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим учителя
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TeacherPersonalAccountHandler - Пытаемся найти учителя", "id_user", session.UserID)
	var teacher models.Teacher
	result = config.DB.Where("id_user = ?", session.UserID).First(&teacher)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	data := TeacherPersonalAccountData{Teacher: teacher}

	templates.ExecuteTemplate(w, "teacher_personal_account", data)
	slog.Info("TeacherPersonalAccountHandler - Успешно", "id_user", session.UserID)
}

// Расписание (выбор недели)
func TeacherScheduleHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "teacher_schedule", nil)
}

// Расписание (время проведения занятий выбранной недели)
func TeacherSchedulePartHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим учителя
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TeacherSchedulePartHandler - Пытаемся найти учителя", "id_user", session.UserID)
	var teacher models.Teacher
	result = config.DB.Where("id_user = ?", session.UserID).First(&teacher)
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
			Where("id_teacher = ? AND (? < date_begin AND date_end < ?)", teacher.ID, point1, point2).
			Order("date_begin").
			Find(&lessons)
		if result.Error != nil {
			templates.ExecuteTemplate(w, "error", errorServerSide)
			return
		}
		weekLessons[ind] = lessons
		point1, point2 = point2, point2.Add(24*time.Hour)
	}

	data := TeacherSchedulePartData{WeekLessons: weekLessons}

	templates.ExecuteTemplate(w, "teacher_schedule_part", data)
	slog.Info("TeacherSchedulePartHandler - Успешно", "id_user", session.UserID)
}

// Дисциплины (возврат дисциплин, которые ведёт преподаватель)
func TeacherDisciplinesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим учителя
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TeacherDisciplinesHandler - Пытаемся найти учителя", "id_user", session.UserID)
	var teacher models.Teacher
	result = config.DB.Where("id_user = ?", session.UserID).First(&teacher)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	data := TeacherDisciplinesData{}

	result = config.DB.
		Table("group_disciplines").
		Where("id_teacher = ?", teacher.ID).
		Group("id_discipline").
		Select(`
        MIN(id) as id,
        MIN(id_group) as id_group,
        id_teacher,
        id_discipline
    `).
		Preload("Discipline").
		Find(&data.GroupDisciplines)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	templates.ExecuteTemplate(w, "teacher_disciplines", data)
	slog.Info("TeacherDisciplinesHandler - Успешно", "id_user", session.UserID)
}

// Дисциплины (возврат групп выбранной дисциплины)
func TeacherDisciplinesPartGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим учителя
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TeacherDisciplinesPartGroupHandler - Пытаемся найти учителя", "id_user", session.UserID)
	var teacher models.Teacher
	result = config.DB.Where("id_user = ?", session.UserID).First(&teacher)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	id_disc, _ := strconv.Atoi(r.FormValue("id_discipline"))

	data := TeacherDisciplinesPartGroupData{}
	result = config.DB.
		Preload("Group").
		Preload("Discipline").
		Where("id_teacher = ? AND id_discipline = ?", teacher.ID, id_disc).
		Find(&data.GroupDisciplines)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	templates.ExecuteTemplate(w, "teacher_disciplines_part_group", data)
	slog.Info("TeacherDisciplinesPartGroupHandler - Успешно", "id_user", session.UserID)

}

// Дисциплины (возврат журнала выбранной группы)
func TeacherDisciplinesPartTableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим учителя
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TeacherDisciplinesPartTableHandler - Пытаемся найти учителя", "id_user", session.UserID)
	var teacher models.Teacher
	result = config.DB.Where("id_user = ?", session.UserID).First(&teacher)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	id_group, _ := strconv.Atoi(r.FormValue("id_group"))
	id_discipline, _ := strconv.Atoi(r.FormValue("id_discipline"))

	// добавить проверку на возможность редактирования (никто другой не редактирует)

	data := TeacherDisciplinesPartTableData{}

	result = config.DB.Where("id = ?", id_discipline).First(&data.Discipline)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	result = config.DB.Where("id = ?", id_group).First(&data.Group)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	result = config.DB.
		Where("id_group = ?", id_group).
		Find(&data.Students)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	result = config.DB.
		Where("id_discipline = ?", id_discipline).
		Find(&data.Lessons)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	data.Actions = make([][]models.Action, len(data.Students))
	for i, student := range data.Students {
		row := make([]models.Action, len(data.Lessons))
		for j, lesson := range data.Lessons {
			result = config.DB.
				Where("id_student = ? AND id_lesson = ?", student.ID, lesson.ID).
				Find(&row[j])
			if result.Error != nil {
				templates.ExecuteTemplate(w, "error", errorServerSide)
				return
			}
		}
		data.Actions[i] = row
	}

	templates.ExecuteTemplate(w, "teacher_disciplines_part_table", data)
	slog.Info("TeacherDisciplinesPartTableHandler - Успешно", "id_user", session.UserID)
}

// Cохранение журнала
type UpdateJournalRequest struct {
	GroupID       int `json:"id_group"`
	DisciplineID  int `json:"id_discipline"`
	ActionChanges []struct {
		StudentID  int `json:"student_id"`
		LessonID   int `json:"lesson_id"`
		Grade      int `json:"grade"`
		Attendance int `json:"attendance"`
	} `json:"actionChanges"`
	LessonChanges []struct {
		LessonID    int    `json:"lesson_id"`
		Description string `json:"description"`
	} `json:"lessonChanges"`
}

func UpdateJournalHandler(w http.ResponseWriter, r *http.Request) {

	var req UpdateJournalRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		fmt.Println("error")
		return
	}

	// Сохраняем изменения в ячейках (оценки/посещаемость)
	for _, change := range req.ActionChanges {

		var action models.Action

		result := config.DB.Where(
			"id_student = ? AND id_lesson = ?",
			change.StudentID,
			change.LessonID,
		).First(&action)

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			action = models.Action{
				StudentID: change.StudentID,
				LessonID:  change.LessonID,
			}
		} else if result.Error != nil {
			slog.Error("DB error", "error", result.Error)
			continue
		}

		// ОЦЕНКА
		if change.Grade != 0 {
			action.Grade = change.Grade
		}

		// ПОСЕЩАЕМОСТЬ
		if change.Attendance != 0 {
			action.Attendance = models.AttendanceType(change.Attendance)
		}

		// ОЧИСТКА (если оба 0)
		if change.Grade == 0 && change.Attendance == 0 {
			action.Grade = 0
			action.Attendance = 0
		}

		config.DB.Save(&action)
	}

	// Сохраняем изменения в уроках (описание)
	for _, change := range req.LessonChanges {
		result := config.DB.Model(&models.Lesson{}).
			Where("id = ?", change.LessonID).
			Update("description", change.Description)

		if result.Error != nil {
			slog.Error("Failed to update lesson description",
				"lesson_id", change.LessonID,
				"error", result.Error)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func mapAttendance(attendance string) int {
	switch attendance {
	case "Я":
		return 0
	case "Н":
		return 1
	case "Б":
		return 2
	case "ДО":
		return 3
	default:
		return 0
	}
}
