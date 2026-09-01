package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthCheck handles checking the health status of the API server and database connection.
// @Summary Check API Health
// @Description Health check endpoint to verify server and database connectivity status
// @Tags Health
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 503 {object} map[string]string
// @Router /health [get]
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	dbStatus := "connected"
	if h.db != nil {
		sqlDB, err := h.db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "disconnected"
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":   "ERROR",
				"message":  "GoTicket API is running, but database connection failed",
				"database": dbStatus,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "OK",
		"message":  "GoTicket API and Database are running cleanly",
		"database": dbStatus,
	})
}
