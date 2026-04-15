package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Date        time.Time      `json:"date"`
	Description string         `json:"description"`
	Total       float64        `json:"total"`
	Type        string         `json:"type"` // "pengeluaran", "pemasukan"
	Attachment  string         `json:"attachment"`

	CategoryID uint     `json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`

	PersonID uint   `json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID" json:"-"`
}
