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
	Name      string         `json:"name"`
	Username  string         `json:"username"`
	Password  string         `json:"password"`
	IsActive  bool           `gorm:"default:false" json:"is_active"`

	Categories   []Category    `gorm:"foreignKey:PersonID" json:"categories"`
	Transactions []Transaction `gorm:"foreignKey:PersonID" json:"transactions"`
}
