package models

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Name      string         `json:"name"`
	Status    string         `json:"status" gorm:"default:'active'"`
	Items     []CartItem     `json:"items" gorm:"foreignKey:CartID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type CartItem struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	CartID   uint    `json:"cart_id" gorm:"not null"`
	ItemID   uint    `json:"item_id" gorm:"not null"`
	Item     Item    `json:"item" gorm:"foreignKey:ItemID"`
	Quantity int     `json:"quantity" gorm:"not null;default:1"`
	Price    float64 `json:"price" gorm:"not null"`
} 