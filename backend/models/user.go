package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"unique;not null"`
	Password  string         `json:"password" gorm:"not null"`
	Token     string         `json:"token"`
	CartID    *uint          `json:"cart_id" gorm:"unique"`
	Cart      *Cart          `json:"cart" gorm:"foreignKey:CartID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
} 