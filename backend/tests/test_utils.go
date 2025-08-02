package tests

import (
	"shopping-cart/models"
	"shopping-cart/routes"
	"shopping-cart/utils"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"gorm.io/gorm"
)

var testDB *gorm.DB

// SetupTestDB initializes an in-memory SQLite database for testing
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	// Auto migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.Item{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
	)
	if err != nil {
		panic("Failed to migrate test database: " + err.Error())
	}

	return db
}

// SetupTestRouter creates a test router with the test database
func SetupTestRouter() *gin.Engine {
	// Setup test database
	testDB = SetupTestDB()
	utils.DB = testDB
	SeedTestData(testDB)

	// Setup router
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routes.SetupRoutes(router)
	
	return router
}

// CleanupTestDB cleans up the test database
func CleanupTestDB(db *gorm.DB) {
	// Delete all data from tables
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM cart_items")
	db.Exec("DELETE FROM carts")
	db.Exec("DELETE FROM items")
	db.Exec("DELETE FROM users")
}

// SeedTestData adds some test data to the database
func SeedTestData(db *gorm.DB) {
	// Create test items
	items := []models.Item{
		{
			Name:        "Test Item 1",
			Description: "First test item",
			Price:       10.99,
			Category:    "Electronics",
			Rating:      4.5,
			Reviews:     10,
			Image:       "https://example.com/item1.jpg",
			InStock:     true,
		},
		{
			Name:        "Test Item 2",
			Description: "Second test item",
			Price:       20.99,
			Category:    "Books",
			Rating:      4.0,
			Reviews:     5,
			Image:       "https://example.com/item2.jpg",
			InStock:     true,
		},
	}

	for _, item := range items {
		db.Create(&item)
	}
}

// CreateTestUser creates a test user and returns the user object
func CreateTestUser(db *gorm.DB, username, password string) models.User {
	user := models.User{
		Username: username,
		Password: password, // In real app, this would be hashed
	}
	db.Create(&user)
	return user
}

// CreateTestCart creates a test cart for a user
func CreateTestCart(db *gorm.DB, userID uint) models.Cart {
	cart := models.Cart{
		UserID: userID,
		Status: "active",
	}
	db.Create(&cart)
	return cart
}

// AddItemToCart adds an item to a cart
func AddItemToCart(db *gorm.DB, cartID, itemID uint, quantity int, price float64) models.CartItem {
	cartItem := models.CartItem{
		CartID:   cartID,
		ItemID:   itemID,
		Quantity: quantity,
		Price:    price,
	}
	db.Create(&cartItem)
	return cartItem
}

// TestSuiteSetup sets up the test suite
func TestSuiteSetup() {
	ginkgo.BeforeEach(func() {
		// Setup test database
		testDB = SetupTestDB()
		
		// Replace the main DB with test DB
		utils.DB = testDB
		
		// Seed test data
		SeedTestData(testDB)
	})

	ginkgo.AfterEach(func() {
		// Cleanup test database
		CleanupTestDB(testDB)
	})
}

// AssertUserExists checks if a user exists in the database
func AssertUserExists(t *testing.T, db *gorm.DB, username string) {
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(user.Username).To(gomega.Equal(username))
}

// AssertUserNotExists checks if a user does not exist in the database
func AssertUserNotExists(t *testing.T, db *gorm.DB, username string) {
	var user models.User
	err := db.Where("username = ?", username).First(&user).Error
	gomega.Expect(err).NotTo(gomega.BeNil())
}

// AssertCartExists checks if a cart exists for a user
func AssertCartExists(t *testing.T, db *gorm.DB, userID uint) {
	var cart models.Cart
	err := db.Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(cart.UserID).To(gomega.Equal(userID))
}

// AssertOrderExists checks if an order exists for a user
func AssertOrderExists(t *testing.T, db *gorm.DB, userID uint) {
	var order models.Order
	err := db.Where("user_id = ?", userID).First(&order).Error
	gomega.Expect(err).To(gomega.BeNil())
	gomega.Expect(order.UserID).To(gomega.Equal(userID))
} 