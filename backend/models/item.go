package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	Price       float64        `json:"price" gorm:"not null"`
	Category    string         `json:"category"`
	Rating      float64        `json:"rating"`
	Reviews     int            `json:"reviews"`
	Image       string         `json:"image"`
	InStock     bool           `json:"in_stock" gorm:"default:true"`
	Status      string         `json:"status" gorm:"default:'active'"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
} 