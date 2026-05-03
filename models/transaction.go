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
	Date        time.Time      `json:"date" validate:"required"`
	Description string         `json:"description" validate:"required"`
	Total       float64        `json:"total" validate:"gt=0"`
	Type        string         `json:"type" validate:"required,oneof=pemasukan pengeluaran"` // "pengeluaran", "pemasukan"
	Attachment  string         `json:"attachment"`

	CategoryID uint     `json:"category_id" validate:"required"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category" validate:"-"`

	PersonID uint   `json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID" json:"-" validate:"-"`
}
