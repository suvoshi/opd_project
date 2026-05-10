package handlers

import (
	"log/slog"
	"net/http"
	"opd_project/config"
	"opd_project/models"
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

type TutorDisciplinesData struct {
	Tables []models.GroupDisciplineTable
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
