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

// signupTestUserCart is declared only once at the top of the file
func signupTestUserCart(router *gin.Engine, username, password string) {
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

func TestCartAPI(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	// Create test user via signup endpoint
	signupTestUserCart(router, "carttest", "password123")
	
	// Login to get token
	loginData := map[string]interface{}{
		"username": "carttest",
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

	t.Run("should add item to cart with valid token", func(t *testing.T) {
		// Prepare cart request
		cartData := map[string]interface{}{
			"item_id":  1,
			"quantity": 2,
		}
		jsonData, _ := json.Marshal(cartData)

		// Make request with token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Item added to cart successfully", response["message"])

		// Verify cart exists in database
		var user models.User
		err = testDB.Where("username = ?", "carttest").First(&user).Error
		assert.NoError(t, err)
		var cart models.Cart
		err = testDB.Where("user_id = ? AND status = ?", user.ID, "active").First(&cart).Error
		assert.NoError(t, err)
		assert.Equal(t, user.ID, cart.UserID)

		// Verify cart item exists
		var cartItem models.CartItem
		err = testDB.Where("cart_id = ? AND item_id = ?", cart.ID, 1).First(&cartItem).Error
		assert.NoError(t, err)
		assert.Equal(t, 2, cartItem.Quantity)
	})

	t.Run("should reject adding item without token", func(t *testing.T) {
		// Prepare cart request
		cartData := map[string]interface{}{
			"item_id":  1,
			"quantity": 2,
		}
		jsonData, _ := json.Marshal(cartData)

		// Make request without token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject adding item with invalid token", func(t *testing.T) {
		// Prepare cart request
		cartData := map[string]interface{}{
			"item_id":  1,
			"quantity": 2,
		}
		jsonData, _ := json.Marshal(cartData)

		// Make request with invalid token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer invalid-token")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("should reject adding non-existent item", func(t *testing.T) {
		// Prepare cart request with non-existent item
		cartData := map[string]interface{}{
			"item_id":  999,
			"quantity": 2,
		}
		jsonData, _ := json.Marshal(cartData)

		// Make request with token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/carts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusNotFound, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Item not found")
	})
}

func TestGetCart(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	// Create test user via signup endpoint
	signupTestUserCart(router, "carttest", "password123")
	
	// Login to get token
	loginData := map[string]interface{}{
		"username": "carttest",
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

	// Add an item to cart first
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

	t.Run("should get user's cart with valid token", func(t *testing.T) {
		// Make request with token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/carts", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["cart"])

		// Verify cart data
		cart := response["cart"].(map[string]interface{})
		var user models.User
		err = testDB.Where("username = ?", "carttest").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, float64(user.ID), cart["user_id"])
		assert.NotNil(t, cart["items"])
	})

	t.Run("should reject getting cart without token", func(t *testing.T) {
		// Make request without token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/carts", nil)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestRemoveFromCart(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	// Create test user via signup endpoint
	signupTestUserCart(router, "carttest", "password123")
	
	// Login to get token
	loginData := map[string]interface{}{
		"username": "carttest",
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

	// Add an item to cart first
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

	t.Run("should remove item from cart with valid token", func(t *testing.T) {
		// Prepare remove request
		removeData := map[string]interface{}{
			"item_id": 1,
		}
		jsonData, _ := json.Marshal(removeData)

		// Make request with token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/carts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Item removed from cart successfully", response["message"])

		// Verify item is removed from cart
		var cartItem models.CartItem
		err = testDB.Where("item_id = ?", 1).First(&cartItem).Error
		assert.Error(t, err) // Should not find the item
	})

	t.Run("should reject removing item without token", func(t *testing.T) {
		// Prepare remove request
		removeData := map[string]interface{}{
			"item_id": 1,
		}
		jsonData, _ := json.Marshal(removeData)

		// Make request without token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("DELETE", "/carts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
} 