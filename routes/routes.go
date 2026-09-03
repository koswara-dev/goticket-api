package routes

import (
	"gotiket-api/handler"
	"gotiket-api/middleware"
	"gotiket-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterConfig struct {
	AppMode               string
	HealthHandler         *handler.HealthHandler
	CustomerHandler       *handler.CustomerHandler
	TicketCategoryHandler *handler.TicketCategoryHandler
	BookingHandler        *handler.BookingHandler
	AuthHandler           *handler.AuthHandler
	ConcertHandler        *handler.ConcertHandler
	AuditLogHandler       *handler.AuditLogHandler
	BlacklistedTokenRepo  repository.BlacklistedTokenRepository
}

func SetupRoutes(r *gin.Engine, cfg RouterConfig) {
	// Register Terpusat Structured JSON Logger Middleware
	r.Use(middleware.GinLogger())

	// Health Check Route
	r.GET("/health", cfg.HealthHandler.HealthCheck)

	// Swagger UI Route (Disabled / Returns 404 in Production Mode)
	if cfg.AppMode == "production" || cfg.AppMode == "prod" {
		r.GET("/swagger/*any", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Resource tidak ditemukan",
			})
		})
	} else {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API v1 Routes
	api := r.Group("/api/v1")
	{
		// Public routes
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", cfg.AuthHandler.Register)
			authGroup.POST("/login", cfg.AuthHandler.Login)
			authGroup.POST("/verify-otp", cfg.AuthHandler.VerifyOTP)
			authGroup.POST("/resend-otp", cfg.AuthHandler.ResendOTP)
			authGroup.POST("/forgot-password", cfg.AuthHandler.ForgotPassword)
			authGroup.POST("/reset-password", cfg.AuthHandler.ResetPassword)
		}

		api.GET("/ticket-categories", cfg.TicketCategoryHandler.GetCategories)
		api.GET("/concerts", cfg.ConcertHandler.FindAll)
		api.GET("/concerts/:id", cfg.ConcertHandler.FindByID)

		// Protected routes (Perlu login JWT)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.BlacklistedTokenRepo))
		{
			// Protected Auth
			protected.POST("/auth/logout", cfg.AuthHandler.Logout)
			protected.GET("/auth/me", cfg.AuthHandler.Me)

			// Protected Customer
			protected.GET("/customers/:id", cfg.CustomerHandler.FindByID)
			protected.POST("/customers", cfg.CustomerHandler.Create)
			protected.PUT("/customers/:id", cfg.CustomerHandler.Update)

			// Protected Booking
			protected.POST("/bookings", cfg.BookingHandler.CreateBooking)
			protected.GET("/bookings/:id", cfg.BookingHandler.GetBookingByID)
		}

		// Admin Only routes (Role = admin)
		adminOnly := protected.Group("")
		adminOnly.Use(middleware.RequiredRole("admin"))
		{
			adminOnly.GET("/customers", cfg.CustomerHandler.FindAll)
			adminOnly.DELETE("/customers/:id", cfg.CustomerHandler.Delete)
			adminOnly.POST("/concerts", cfg.ConcertHandler.Create)
			adminOnly.PUT("/concerts/:id", cfg.ConcertHandler.Update)
			adminOnly.DELETE("/concerts/:id", cfg.ConcertHandler.Delete)
			adminOnly.GET("/bookings/report/pdf", cfg.BookingHandler.ExportBookingReportPDF)
			adminOnly.GET("/audit-logs", cfg.AuditLogHandler.GetAuditLogs)
		}
	}
}
