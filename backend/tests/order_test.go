package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shopping-cart/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// signupTestUserOrder is declared only once at the top of the file
func signupTestUserOrder(router *gin.Engine, username, password string) {
	userData := map[string]interface{}{
		"username": username,
		"password": password,
	}
	jsonData, _ := json.Marshal(userData)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
}

func TestCreateOrder(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	// Create test user via signup endpoint
	signupTestUserOrder(router, "ordertest", "password123")
	
	// Login to get token
	loginData := map[string]interface{}{
		"username": "ordertest",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	token := response["token"].(string)

	// Add items to cart first
	cartData1 := map[string]interface{}{
		"item_id":  1,
		"quantity": 2,
	}
	jsonData1, _ := json.Marshal(cartData1)

	w = httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonData1))
	req1.Header.Set("Content-Type", "application/json")
	req1.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req1)

	// Add second item
	cartData2 := map[string]interface{}{
		"item_id":  2,
		"quantity": 1,
	}
	jsonData2, _ := json.Marshal(cartData2)

	w = httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonData2))
	req2.Header.Set("Content-Type", "application/json")
	req2.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req2)

	t.Run("should create order from cart with valid token", func(t *testing.T) {
		// Make order request with token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/orders", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Order created successfully", response["message"])

		// Verify order exists in database
		var user models.User
		err = testDB.Where("username = ?", "ordertest").First(&user).Error
		assert.NoError(t, err)
		var order models.Order
		err = testDB.Where("user_id = ?", user.ID).First(&order).Error
		assert.NoError(t, err)
		assert.Equal(t, user.ID, order.UserID)
		assert.Greater(t, order.Total, 0.0)

		// Verify cart is cleared (no active cart items)
		var cartItems []models.CartItem
		err = testDB.Where("cart_id = ?", order.CartID).Find(&cartItems).Error
		assert.NoError(t, err)
		assert.Equal(t, 0, len(cartItems))
	})

	t.Run("should reject creating order without token", func(t *testing.T) {
		// Make order request without token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/orders", nil)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject creating order with invalid token", func(t *testing.T) {
		// Make order request with invalid token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/orders", nil)
		req.Header.Set("Authorization", "Bearer invalid-token")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject creating order with empty cart", func(t *testing.T) {
		// Create a new user with empty cart
		signupTestUserOrder(router, "emptycart", "password123")
		
		// Login to get token for empty user
		loginData := map[string]interface{}{
			"username": "emptycart",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		emptyToken := response["token"].(string)

		// Try to create order with empty cart
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("POST", "/orders", nil)
		req.Header.Set("Authorization", "Bearer "+emptyToken)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusNotFound, w.Code)

		var orderResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &orderResponse)
		assert.NoError(t, err)
		assert.Contains(t, orderResponse["error"], "Cart not found")
	})
}

func TestListOrders(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	// Create test user via signup endpoint
	signupTestUserOrder(router, "ordertest", "password123")
	
	// Login to get token
	loginData := map[string]interface{}{
		"username": "ordertest",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	token := response["token"].(string)

	// Add items to cart and create order
	cartData := map[string]interface{}{
		"item_id":  1,
		"quantity": 2,
	}
	jsonData, _ = json.Marshal(cartData)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// Create order
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/orders", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	t.Run("should list user's orders with valid token", func(t *testing.T) {
		// Make request with token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["orders"])

		// Verify orders data
		orders := response["orders"].([]interface{})
		assert.Equal(t, 1, len(orders))

		orderData := orders[0].(map[string]interface{})
		var user models.User
		err = testDB.Where("username = ?", "ordertest").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, float64(user.ID), orderData["user_id"])
		assert.Greater(t, orderData["total"], 0.0)
	})

	t.Run("should reject listing orders without token", func(t *testing.T) {
		// Make request without token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/orders", nil)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestOrderDataIntegrity(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	// Create test user via signup endpoint
	signupTestUserOrder(router, "ordertest", "password123")
	
	// Login to get token
	loginData := map[string]interface{}{
		"username": "ordertest",
		"password": "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	token := response["token"].(string)

	t.Run("should maintain order data integrity", func(t *testing.T) {
		// Add items to cart
		cartData1 := map[string]interface{}{
			"item_id":  1,
			"quantity": 2,
		}
		jsonData1, _ := json.Marshal(cartData1)

		w := httptest.NewRecorder()
		req1, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")
		req1.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req1)

		// Create order
		w = httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/orders", nil)
		req2.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req2)

		// Verify order exists
		var user models.User
		err := testDB.Where("username = ?", "ordertest").First(&user).Error
		assert.NoError(t, err)
		var order models.Order
		err = testDB.Where("user_id = ?", user.ID).First(&order).Error
		assert.NoError(t, err)

		// Verify cart is linked to order
		assert.NotEqual(t, uint(0), order.CartID)

		// Verify cart exists and is linked
		var cart models.Cart
		err = testDB.Where("id = ?", order.CartID).First(&cart).Error
		assert.NoError(t, err)
		assert.Equal(t, user.ID, cart.UserID)

		// Verify order total is calculated correctly
		expectedTotal := 10.99 * 2 // item price * quantity
		assert.Equal(t, expectedTotal, order.Total)
	})
} 