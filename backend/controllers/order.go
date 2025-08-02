package controllers

import (
	"net/http"
	"shopping-cart/models"
	"shopping-cart/utils"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	userID := c.GetUint("user_id")

	// Get user's active cart
	var cart models.Cart
	if err := utils.DB.Where("user_id = ? AND status = ?", userID, "active").Preload("Items.Item").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	if len(cart.Items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Calculate total
	var total float64
	for _, item := range cart.Items {
		total += item.Price * float64(item.Quantity)
	}

	// Create order from cart (as per ERD)
	order := models.Order{
		CartID: cart.ID,
		UserID: userID,
		Total:  total,
		Status: "pending",
	}

	if err := utils.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Clear cart items after order creation
	utils.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order_id": order.ID,
		"total":    total,
	})
}

func ListOrders(c *gin.Context) {
	userID := c.GetUint("user_id")

	var orders []models.Order
	if err := utils.DB.Where("user_id = ?", userID).Preload("Cart.Items.Item").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func ListAllOrders(c *gin.Context) {
	var orders []models.Order
	if err := utils.DB.Preload("User").Preload("Cart.Items.Item").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
} 