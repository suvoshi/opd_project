package main

import (
	"io"
	"log/slog"
	"net/http"
	"opd_project/config"
	"opd_project/handlers"
	"opd_project/models"
	"os"
)

func main() {
	// 0. Настройка логирования
	file, err := os.OpenFile("logging.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)
	handler := slog.New(slog.NewJSONHandler(multiWriter, nil))
	slog.SetDefault(handler)

	// 1. Инициализация БД и таблиц
	config.InitDB()
	config.DB.AutoMigrate(
		&models.User{},
		&models.Student{},
		&models.Lesson{},
		&models.Group{},
		&models.Discipline{},
		&models.GroupDiscipline{},
		&models.Action{},
		&models.StudentEndDiscipline{},
		&models.Session{},
		&models.Teacher{},
		&models.Tutor{},
		&models.Announcement{},
	)

	// 2. Загрузка шаблонов
	handlers.InitTemplates()

	// 3. Настройка роутов (Маршрутизация)
	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login/", handlers.LoginHandler)
	http.HandleFunc("/student/", handlers.StudentHandler)
	http.HandleFunc("/teacher/", handlers.TeacherHandler)
	http.HandleFunc("/tutor/", handlers.TutorHandler)

	// 3. 1 API для HTMX
	http.HandleFunc("/try_login", handlers.TryLogin)

	// 3.1.1 Для студента
	http.HandleFunc("/student/dashboard/", handlers.StudentDashboardHandler)
	http.HandleFunc("/student/personal_account/", handlers.StudentPersonalAccountHandler)
	http.HandleFunc("/student/schedule/", handlers.StudentScheduleHandler)
	http.HandleFunc("/student/schedule/part/", handlers.StudentSchedulePartHandler)
	http.HandleFunc("/student/discipline_progress/", handlers.StudentDisciplineProgressHandler)

	// 4. Раздача статики (CSS, JS, картинки)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 5. Запуск сервера
	slog.Info("Сервер запущен на http://localhost:8080")
	err = http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		slog.Error("Ошибка сервера")
	}
}
