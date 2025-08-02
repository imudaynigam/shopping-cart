package utils

import (
	"log"
	"os"
	"shopping-cart/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	
	// Get database configuration from environment
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "sqlite"
	}
	
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "shopping_cart.db"
	}
	
	// Initialize database based on type
	if dbType == "sqlite" {
		DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	} else {
		// For future PostgreSQL support
		log.Fatal("Only SQLite is currently supported")
	}
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed initial data
	seedData()
}

func seedData() {
	// Check if items already exist
	var count int64
	DB.Model(&models.Item{}).Count(&count)
	if count > 0 {
		return // Data already seeded
	}

	// Seed items
	items := []models.Item{
		{
			Name:        "Premium Wireless Headphones",
			Description: "High-quality wireless headphones with active noise cancellation and 30-hour battery life",
			Price:       299.99,
			Category:    "Electronics",
			Rating:      4.8,
			Reviews:     1247,
			Image:       "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=400&h=300&fit=crop",
			InStock:     true,
		},
		{
			Name:        "Smart Fitness Watch",
			Description: "Advanced fitness tracking with heart rate monitoring, GPS, and 7-day battery life",
			Price:       199.99,
			Category:    "Electronics",
			Rating:      4.6,
			Reviews:     892,
			Image:       "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=400&h=300&fit=crop",
			InStock:     true,
		},
		{
			Name:        "Ergonomic Office Chair",
			Description: "Premium ergonomic office chair with adjustable lumbar support and memory foam cushion",
			Price:       449.99,
			Category:    "Furniture",
			Rating:      4.7,
			Reviews:     456,
			Image:       "https://images.unsplash.com/photo-1567538096630-e0c55bd6374c?w=400&h=300&fit=crop",
			InStock:     true,
		},
		{
			Name:        "Organic Coffee Beans",
			Description: "Premium organic coffee beans from sustainable farms in Colombia",
			Price:       24.99,
			Category:    "Food & Beverages",
			Rating:      4.9,
			Reviews:     2341,
			Image:       "https://images.unsplash.com/photo-1559056199-641a0ac8b55e?w=400&h=300&fit=crop",
			InStock:     true,
		},
		{
			Name:        "Professional Camera Lens",
			Description: "85mm f/1.4 portrait lens with beautiful bokeh and exceptional sharpness",
			Price:       899.99,
			Category:    "Electronics",
			Rating:      4.9,
			Reviews:     567,
			Image:       "https://images.unsplash.com/photo-1516035069371-29a1b244cc32?w=400&h=300&fit=crop",
			InStock:     true,
		},
	}

	for _, item := range items {
		DB.Create(&item)
	}

	log.Println("Database seeded successfully")
} 