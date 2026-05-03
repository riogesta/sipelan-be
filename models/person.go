package models

import (
	"time"

	"gorm.io/gorm"
)

type Person struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name      string         `json:"name" validate:"required"`
	Username  string         `json:"username" validate:"required"`
	Password  string         `json:"password" validate:"required,min=6"`
	IsActive  bool           `gorm:"default:false" json:"is_active"`

	Categories   []Category    `gorm:"foreignKey:PersonID" json:"categories"`
	Transactions []Transaction `gorm:"foreignKey:PersonID" json:"transactions"`
}
