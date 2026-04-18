package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"opd_project/config"
	"opd_project/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	//"strconv"
)

// HTMX хендлер для логина
func TryLogin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("HX-Request") != "true" {
		http.Error(w, "This endpoint requires HTMX request", http.StatusForbidden)
		return
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
