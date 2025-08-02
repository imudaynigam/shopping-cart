package routes

import (
	"net/http"
	"shopping-cart/controllers"
	"shopping-cart/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Root route to show server is running
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Shopping Cart Backend Server is running!",
			"version": "1.0.0",
			"status":  "active",
			"endpoints": gin.H{
				"auth": gin.H{
					"POST /users": "Sign up a new user",
					"POST /users/login": "Login user",
					"GET /users": "List all users (protected)",
				},
				"items": gin.H{
					"GET /items": "List all items",
					"POST /items": "Create new item (protected)",
					"DELETE /items/:id": "Delete item (protected)",
				},
				"cart": gin.H{
					"POST /carts": "Add item to cart (protected)",
					"DELETE /carts": "Remove item from cart (protected)",
					"GET /carts": "Get user's cart (protected)",
					"GET /carts/all": "List all carts (protected)",
				},
				"orders": gin.H{
					"POST /orders": "Create order from cart (protected)",
					"GET /orders": "List user's orders (protected)",
					"GET /orders/all": "List all orders (protected)",
				},
			},
		})
	})

	// Public routes
	public := r.Group("/")
	{
		public.POST("/users", controllers.Signup)
		public.POST("/users/login", controllers.Login)
		public.GET("/items", controllers.ListItems)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middlewares.AuthMiddleware())
	{
		// User routes
		protected.GET("/users", controllers.ListUsers)

		// Item routes
		protected.POST("/items", controllers.CreateItem)
		protected.DELETE("/items/:id", controllers.DeleteItem)

		// Cart routes
		protected.POST("/carts", controllers.AddToCart)
		protected.DELETE("/carts", controllers.RemoveFromCart)
		protected.GET("/carts", controllers.GetCart)
		protected.GET("/carts/all", controllers.ListCarts) // Admin only

		// Order routes
		protected.POST("/orders", controllers.CreateOrder)
		protected.GET("/orders", controllers.ListOrders)
		protected.GET("/orders/all", controllers.ListAllOrders) // Admin only
	}
} 