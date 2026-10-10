package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/toko-berkah/pos-backend/handlers"
	"github.com/toko-berkah/pos-backend/middleware"
	"github.com/toko-berkah/pos-backend/repositories"
	"github.com/toko-berkah/pos-backend/services"
	"github.com/toko-berkah/pos-backend/utils"
	"github.com/toko-berkah/pos-backend/config"
)

// SetupRoutes initializes all API routes
func SetupRoutes(db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()

	// Load config for JWT
	cfg := config.LoadConfig()
	utils.InitJWT(cfg)

	// Apply global middleware
	router.Use(middleware.ErrorHandlerMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.FrontendURL))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	productRepo := repositories.NewProductRepository(db)
	inventoryRepo := repositories.NewInventoryRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	productService := services.NewProductService(productRepo, inventoryRepo)
	inventoryService := services.NewInventoryService(inventoryRepo, productRepo)
	transactionService := services.NewTransactionService(db, productRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	productHandler := handlers.NewProductHandler(productService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			// Note: Public /register endpoint removed for MVP (single operator only)
			auth.GET("/profile", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), authHandler.GetProfile)
			auth.POST("/logout", middleware.AuthMiddleware(db), authHandler.Logout)
			auth.POST("/change-password", middleware.AuthMiddleware(db), authHandler.ChangePassword)
		}

		// Product routes
		products := v1.Group("/products")
		{
			// Public endpoints
			products.GET("", productHandler.GetAllProducts)
			products.GET("/sku", productHandler.GetProductBySKU)
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/:id", productHandler.GetProduct)

			// Protected endpoints (authentication required, no role check for MVP)
			products.POST("", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), productHandler.CreateProduct)
			products.PUT("/:id", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), productHandler.UpdateProduct)
		}

		// Inventory routes
		inventory := v1.Group("/inventory")
		{
			// Protected endpoints
			inventory.GET("", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), inventoryHandler.GetInventory)
			inventory.GET("/all", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), inventoryHandler.GetAllInventory)
			inventory.GET("/low-stock", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), inventoryHandler.GetLowStockProducts)
			inventory.POST("/adjustment", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), inventoryHandler.AdjustInventory)
			inventory.GET("/logs", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), inventoryHandler.GetInventoryLogs)
			inventory.GET("/logs/product", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), inventoryHandler.GetInventoryLogsByProduct)
		}

		// Transaction routes
		transactions := v1.Group("/transactions")
		{
			// Protected endpoints
			transactions.POST("", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), transactionHandler.CreateTransaction)
			transactions.GET("", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), transactionHandler.GetAllTransactions)
			transactions.GET("/:id", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), transactionHandler.GetTransaction)
			transactions.POST("/:id/confirm-payment", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), transactionHandler.ConfirmPayment)
			transactions.POST("/:id/cancel", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), transactionHandler.CancelTransaction)
			transactions.GET("/date-range", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), transactionHandler.GetTransactionsByDateRange)
		}

		// Report routes (placeholder for next priority)
		reports := v1.Group("/reports")
		{
			reports.GET("/dashboard", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "dashboard endpoint - coming soon"})
			})
			reports.GET("/sales-breakdown", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "sales breakdown endpoint - coming soon"})
			})
			reports.GET("/bestseller", middleware.AuthMiddleware(db), middleware.RequirePasswordChange(), func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "bestseller endpoint - coming soon"})
			})
		}
	}

	return router
}
