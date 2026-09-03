package main

import (
	"fmt"

	"gotiket-api/config"
	_ "gotiket-api/docs"
	"gotiket-api/handler"
	"gotiket-api/repository"
	"gotiket-api/routes"
	"gotiket-api/service"
	"gotiket-api/utils"

	"github.com/gin-gonic/gin"
)

// @title GoTicket API Documentation
// @version 1.0.1
// @description REST API server for GoTicket ticketing system.
// @termsOfService http://swagger.io/terms/

// @contact.name GoTicket Support Team
// @contact.url http://www.goticket.com/support
// @contact.email support@goticket.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:18080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type 'Bearer ' followed by your JWT token.
func main() {
	// Inisialisasi Structured Logger Terpusat (Rotasi 24 jam & Retensi 365 hari)
	utils.InitLogger()

	appConfig := config.LoadConfig()

	if appConfig.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
		utils.Log.Info("Sistem berjalan dalam Mode Produksi (Gin Release Mode)")
	} else {
		utils.Log.Info("Sistem berjalan dalam Mode Pengembangan (Gin Debug Mode)")
	}

	db := config.InitDB()

	// Storage Provider
	storageProvider := service.NewLocalStorageProvider("./uploads", appConfig.BaseURL)

	// Initialize Repositories
	customerRepo := repository.NewCustomerRepository(db)
	concertRepo := repository.NewConcertRepository(db)
	ticketCategoryRepo := repository.NewTicketCategoryRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	userRepo := repository.NewUserRepository(db)
	blacklistedTokenRepo := repository.NewBlacklistedTokenRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	auditLogRepo := repository.NewAuditLogRepository(db)

	// Initialize Services
	auditLogService := service.NewAuditLogService(auditLogRepo)
	customerService := service.NewCustomerService(customerRepo)
	concertService := service.NewConcertService(concertRepo, storageProvider)
	bookingService := service.NewBookingService(db, bookingRepo, customerRepo)
	authService := service.NewAuthService(userRepo, blacklistedTokenRepo, otpRepo, appConfig)

	// Initialize Handlers
	healthHandler := handler.NewHealthHandler(db)
	customerHandler := handler.NewCustomerHandler(customerService)
	concertHandler := handler.NewConcertHandler(concertService, auditLogService)
	ticketCategoryHandler := handler.NewTicketCategoryHandler(ticketCategoryRepo)
	bookingHandler := handler.NewBookingHandler(bookingService, auditLogService)
	authHandler := handler.NewAuthHandler(authService, auditLogService)
	auditLogHandler := handler.NewAuditLogHandler(auditLogService)

	// Setup Gin Framework Router
	r := gin.Default()

	// Serve static uploads
	r.Static("/uploads", "./uploads")

	// Register Routes
	routes.SetupRoutes(r, routes.RouterConfig{
		AppMode:               appConfig.AppMode,
		HealthHandler:         healthHandler,
		CustomerHandler:       customerHandler,
		ConcertHandler:        concertHandler,
		TicketCategoryHandler: ticketCategoryHandler,
		BookingHandler:        bookingHandler,
		AuthHandler:           authHandler,
		AuditLogHandler:       auditLogHandler,
		BlacklistedTokenRepo:  blacklistedTokenRepo,
	})

	utils.Log.Infof("Server sudah berjalan on port %d...", appConfig.AppPort)
	if err := r.Run(fmt.Sprintf(":%d", appConfig.AppPort)); err != nil {
		utils.Log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

