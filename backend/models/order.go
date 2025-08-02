package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CartID    uint           `json:"cart_id" gorm:"not null"`
	Cart      Cart           `json:"cart" gorm:"foreignKey:CartID"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Total     float64        `json:"total" gorm:"not null"`
	Status    string         `json:"status" gorm:"default:'pending'"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
} 