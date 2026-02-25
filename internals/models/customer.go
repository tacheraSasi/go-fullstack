package models

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"not null" json:"name"`
	Email     string         `gorm:"not null;uniqueIndex" json:"email"`
	Phone     string         `json:"phone"`
	Address   string         `json:"address"`
}
