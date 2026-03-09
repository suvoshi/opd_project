package main

import (
	"net/http"
	"opd_project/config"
	"opd_project/handlers"
	"opd_project/models"
)

func main() {
	// 1. Инициализация БД и таблиц
	config.InitDB()
	config.DB.AutoMigrate(&models.Student{}, &models.Grade{})

	// 2. Загрузка шаблонов
	handlers.InitTemplates()

	// 3. Настройка роутов (Маршрутизация)

	// Главная страница
	http.HandleFunc("/", handlers.IndexHandler)

	// API для HTMX
	http.HandleFunc("/grades", handlers.GetGradesHandler)
	http.HandleFunc("/grade/add", handlers.AddGradeHandler)

	// 4. Раздача статики (CSS, JS, картинки)
	// Все файлы из папки static будут доступны по пути /static/...
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// 5. Запуск сервера
	println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
