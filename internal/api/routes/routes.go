package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/api/handlers"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/api/middleware"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/config"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/security/auth"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/security/crypto"
	"github.com/nshmdayo/in-house-datamanagement-system-sample/internal/services"
)

// SetupRoutes configures all application routes
func SetupRoutes(cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))
	router.Use(middleware.RateLimitMiddleware())
	router.Use(gin.Recovery())

	// Initialize services
	tokenService := auth.NewTokenService(cfg)
	passwordService := crypto.NewPasswordService()
	userService := services.NewUserService()
	auditService := services.NewAuditService()

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(tokenService, passwordService, userService, auditService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "datamanagement-system",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(tokenService, userService))
		{
			// Auth routes
			authProtected := protected.Group("/auth")
			{
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.GET("/profile", authHandler.GetProfile)
			}

			// TODO: Implement additional handlers
			// User management routes
			// users := protected.Group("/users")
			// {
			// 	users.GET("", userHandler.GetUsers)
			// 	users.GET("/:id", userHandler.GetUser)
			// 	users.PUT("/:id", userHandler.UpdateUser)
			// }

			// Document management routes
			// documents := protected.Group("/documents")
			// {
			// 	documents.GET("", documentHandler.GetDocuments)
			// 	documents.POST("", documentHandler.CreateDocument)
			// }

			// Blockchain routes
			// blockchain := protected.Group("/blockchain")
			// {
			// 	blockchain.GET("/blocks", blockchainHandler.GetBlocks)
			// 	blockchain.POST("/verify", blockchainHandler.VerifyIntegrity)
			// }
		}
	}

	return router
}
