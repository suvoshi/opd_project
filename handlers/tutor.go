package handlers

import (
	"log/slog"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"strconv"
	"time"
	//"strconv"
)

// HTMX хендлеры для учителя

// Данные, используемые хендлерами для заполнения шаблонов
type TutorDashboardData struct {
	AnnouncementData []models.Announcement
}

type TutorPersonalAccountData struct {
	Tutor models.Tutor
}

// Дашборд
func TutorDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим куратора
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TutorDashboardHandler - Пытаемся найти куратора", "id_user", session.UserID)
	var tutor models.Tutor
	result = config.DB.Where("id_user = ?", session.UserID).First(&tutor)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	now := time.Now()
	weekAgo := now.Add(-7 * 24 * time.Hour)

	data := TutorDashboardData{}

	result = config.DB.
		Where("(date BETWEEN ? AND ?) AND visibility <= 2", weekAgo, now).
		Order("date DESC").
		Find(&data.AnnouncementData)

	templates.ExecuteTemplate(w, "tutor_dashboard", data)
	slog.Info("TutorDashboardHandler - Успешно", "id_user", session.UserID)
}

// Личный кабинет
func TutorPersonalAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим куратора
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TutorPersonalAccountHandler - Пытаемся найти куратора", "id_user", session.UserID)
	var tutor models.Tutor
	result = config.DB.Where("id_user = ?", session.UserID).First(&tutor)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	data := TutorPersonalAccountData{Tutor: tutor}

	templates.ExecuteTemplate(w, "tutor_personal_account", data)
	slog.Info("TutorPersonalAccountHandler - Успешно", "id_user", session.UserID)
}

// Дисциплины (возврат дисциплин, которые ведёт преподаватель)
func TutorDisciplinesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим куратора
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TutorDisciplinesHandler - Пытаемся найти учителя", "id_user", session.UserID)
	var tutor models.Tutor
	result = config.DB.Where("id_user = ?", session.UserID).First(&tutor)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	data := TeacherDisciplinesData{}

	result = config.DB.
		Table("group_disciplines").
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
	slog.Info("TutorDisciplinesHandler - Успешно", "id_user", session.UserID)
}

// Дисциплины (возврат групп выбранной дисциплины)
func TutorDisciplinesPartGroupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим куратора
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TutorDisciplinesPartGroupHandler - Пытаемся найти куратора", "id_user", session.UserID)
	var tutor models.Tutor
	result = config.DB.Where("id_user = ?", session.UserID).First(&tutor)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	id_disc, _ := strconv.Atoi(r.FormValue("id_discipline"))

	data := TeacherDisciplinesPartGroupData{}
	result = config.DB.
		Preload("Group").
		Preload("Discipline").
		Where("id_discipline = ?", id_disc).
		Find(&data.GroupDisciplines)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}

	templates.ExecuteTemplate(w, "teacher_disciplines_part_group", data)
	slog.Info("TutorDisciplinesPartGroupHandler - Успешно", "id_user", session.UserID)

}

// Дисциплины (возврат журнала выбранной группы)
func TutorDisciplinesPartTableHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// находим куратора
	var session models.Session
	result := config.DB.Where("id_session = ?", cookie.Value).First(&session)
	if result.Error != nil {
		templates.ExecuteTemplate(w, "error", errorServerSide)
		return
	}
	slog.Info("TutorDisciplinesPartTableHandler - Пытаемся найти куратора", "id_user", session.UserID)
	var tutor models.Tutor
	result = config.DB.Where("id_user = ?", session.UserID).First(&tutor)
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
	slog.Info("TutorDisciplinesPartTableHandler - Успешно", "id_user", session.UserID)
}
