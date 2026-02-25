package models

import (
	"time"

	"gorm.io/gorm"
)

type BlacklistedToken struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Token     string         `gorm:"not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
}

type PasswordResetToken struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Token     string         `gorm:"not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expires_at"`
	UsedAt    *time.Time     `json:"used_at,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}
