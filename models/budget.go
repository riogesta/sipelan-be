package models

import "gorm.io/gorm"

type Budget struct {
	gorm.Model
	PersonID   uint     `json:"person_id"`
	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
	Amount     float64  `json:"amount"`
	Month      int      `json:"month"`
	Year       int      `json:"year"`
}
