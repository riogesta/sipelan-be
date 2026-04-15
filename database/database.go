package database

import (
	"log"
	"sipelan/config"
	"sipelan/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.Load()
	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}

	log.Println("Database connection established")

	err = DB.AutoMigrate(&models.Person{}, &models.Category{}, &models.Transaction{})
	if err != nil {
		log.Println("Failed to migrate database:", err)
	}
}
