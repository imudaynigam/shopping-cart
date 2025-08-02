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

func setupTestDB() *gin.Engine {
	return SetupTestRouter()
}

func cleanupTestDB() {
	CleanupTestDB(testDB)
}

// signupTestUser is declared only once at the top of the file
func signupTestUserUser(router *gin.Engine, username, password string) {
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

func TestUserSignup(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	t.Run("should create a new user successfully", func(t *testing.T) {
		// Prepare request
		userData := map[string]interface{}{
			"username": "testuser",
			"password": "testpass123",
		}
		jsonData, _ := json.Marshal(userData)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusCreated, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User created successfully", response["message"])

		// Verify user exists in database
		var user models.User
		err = testDB.Where("username = ?", "testuser").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, "testuser", user.Username)
	})

	t.Run("should reject signup with missing username", func(t *testing.T) {
		// Prepare request without username
		userData := map[string]interface{}{
			"password": "testpass123",
		}
		jsonData, _ := json.Marshal(userData)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Field validation")
	})

	t.Run("should reject signup with duplicate username", func(t *testing.T) {
		// Create first user
		userData1 := map[string]interface{}{
			"username": "duplicateuser",
			"password": "testpass123",
		}
		jsonData1, _ := json.Marshal(userData1)

		w := httptest.NewRecorder()
		req1, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData1))
		req1.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req1)

		assert.Equal(t, http.StatusCreated, w.Code)

		// Try to create second user with same username
		userData2 := map[string]interface{}{
			"username": "duplicateuser",
			"password": "differentpass",
		}
		jsonData2, _ := json.Marshal(userData2)

		w = httptest.NewRecorder()
		req2, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData2))
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req2)

		// Assert response
		assert.Equal(t, http.StatusConflict, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Username already exists")
	})
}

func TestUserLogin(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	// Create a test user via signup endpoint
	signupTestUserUser(router, "logintest", "password123")

	t.Run("should login with valid credentials and return token", func(t *testing.T) {
		// Prepare login request
		loginData := map[string]interface{}{
			"username": "logintest",
			"password": "password123",
		}
		jsonData, _ := json.Marshal(loginData)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Login successful", response["message"])
		assert.NotNil(t, response["token"])
		assert.NotEqual(t, "", response["token"])

		// Verify token is stored in database
		var user models.User
		err = testDB.Where("username = ?", "logintest").First(&user).Error
		assert.NoError(t, err)
		assert.Equal(t, response["token"], user.Token)
	})

	t.Run("should reject login with invalid password", func(t *testing.T) {
		// Prepare login request with wrong password
		loginData := map[string]interface{}{
			"username": "logintest",
			"password": "wrongpassword",
		}
		jsonData, _ := json.Marshal(loginData)

		// Make request
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/users/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid username or password", response["error"])
	})
}

func TestListUsers(t *testing.T) {
	router := setupTestDB()
	defer cleanupTestDB()

	// Create a test user via signup endpoint
	signupTestUserUser(router, "listtest", "password123")
	
	// Login to get token
	loginData := map[string]interface{}{
		"username": "listtest",
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

	t.Run("should list users with valid token", func(t *testing.T) {
		// Make request with token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotNil(t, response["users"])
	})

	t.Run("should reject request without token", func(t *testing.T) {
		// Make request without token
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/users", nil)
		router.ServeHTTP(w, req)

		// Assert response
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
} 