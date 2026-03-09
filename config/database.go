package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("journal.db"), &gorm.Config{})
	if err != nil {
		panic("Не удалось подключиться к базе данных: " + err.Error())
	}
}
