package tests

import (
	"context"
	"os"
	"sipelan/common"
	"sipelan/database"
	"sipelan/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB menyiapkan database SQLite in-memory untuk pengujian isolated
func SetupTestDB() {
	os.Setenv("JWT_SECRET_KEY", "test-secret-key")
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("Gagal menghubungkan ke database test")
	}

	db.AutoMigrate(&models.Person{}, &models.Category{}, &models.Transaction{}, &models.Budget{})
	database.DB = db
}

// CreateTestUser membuat user buatan untuk keperluan otentikasi dalam test
func CreateTestUser() models.Person {
	hashedPassword, _ := common.HashPassword("password123")
	person := models.Person{
		Username: "testuser",
		Password: hashedPassword,
		IsActive: true,
	}
	database.DB.Create(&person)
	return person
}

// GetTestContext menyuntikkan data user ke dalam context request
func GetTestContext(person models.Person) context.Context {
	return context.WithValue(context.Background(), "person", person)
}
