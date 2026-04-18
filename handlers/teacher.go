package handlers

import (
	"log/slog"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"strconv"
	//"strconv"
)

// HTMX хендлеры для учителя

// Данные, используемые хендлерами для заполнения шаблонов
type TeacherDisciplinesData struct {
	GroupDisciplines []models.GroupDiscipline
}

type TeacherDisciplinesPartGroupData struct {
	GroupDisciplines []models.GroupDiscipline
}

// Дашборд
func TeacherDashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "teacher_dashboard", nil)
}

// Личный кабинет
func TeacherPersonalAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
	}
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "teacher_personal_account", nil)
}

// Расписание
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
	_, err := r.Cookie("id_session")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "teacher_disciplines_part_table", nil)
}
