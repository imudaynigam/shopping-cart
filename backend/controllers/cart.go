package controllers

import (
	"net/http"
	"shopping-cart/models"
	"shopping-cart/utils"

	"github.com/gin-gonic/gin"
)

type AddToCartRequest struct {
	ItemID   uint `json:"item_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}

type RemoveFromCartRequest struct {
	ItemID uint `json:"item_id" binding:"required"`
}

func AddToCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if item exists
	var item models.Item
	if err := utils.DB.First(&item, req.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Get user's active cart or create new one
	var cart models.Cart
	if err := utils.DB.Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error; err != nil {
		// Create new cart
		cart = models.Cart{
			UserID: userID,
			Name:   "Shopping Cart",
			Status: "active",
		}
		if err := utils.DB.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
		
		// Update user's cart_id
		utils.DB.Model(&models.User{}).Where("id = ?", userID).Update("cart_id", cart.ID)
	}

	// Check if item already in cart
	var existingCartItem models.CartItem
	if err := utils.DB.Where("cart_id = ? AND item_id = ?", cart.ID, req.ItemID).First(&existingCartItem).Error; err == nil {
		// Update quantity
		existingCartItem.Quantity += req.Quantity
		utils.DB.Save(&existingCartItem)
	} else {
		// Add new item to cart
		cartItem := models.CartItem{
			CartID:   cart.ID,
			ItemID:   req.ItemID,
			Quantity: req.Quantity,
			Price:    item.Price,
		}
		if err := utils.DB.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add item to cart"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

func RemoveFromCart(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req RemoveFromCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user's active cart
	var cart models.Cart
	if err := utils.DB.Where("user_id = ? AND status = ?", userID, "active").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Remove item from cart
	if err := utils.DB.Where("cart_id = ? AND item_id = ?", cart.ID, req.ItemID).Delete(&models.CartItem{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove item from cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}

func GetCart(c *gin.Context) {
	userID := c.GetUint("user_id")

	var cart models.Cart
	if err := utils.DB.Where("user_id = ? AND status = ?", userID, "active").Preload("Items.Item").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func ListCarts(c *gin.Context) {
	var carts []models.Cart
	if err := utils.DB.Preload("User").Preload("Items.Item").Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch carts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"carts": carts})
} 