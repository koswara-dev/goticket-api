package middleware

import (
	"net/http"
	"strings"

	"gotiket-api/repository"
	"gotiket-api/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(blacklistedTokenRepo repository.BlacklistedTokenRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Log.WithFields(map[string]interface{}{
				"client_ip": c.ClientIP(),
				"path":      c.Request.URL.Path,
			}).Warn("Akses ditolak: Authorization header kosong")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Header Authorization dibutuhkan",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.Log.WithFields(map[string]interface{}{
				"client_ip": c.ClientIP(),
				"path":      c.Request.URL.Path,
			}).Warn("Akses ditolak: Format Bearer Token tidak valid")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Format Authorization header harus Bearer <token>",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Cek apakah token di-blacklist
		isBlacklisted, err := blacklistedTokenRepo.IsBlacklisted(tokenString)
		if err != nil || isBlacklisted {
			utils.Log.WithFields(map[string]interface{}{
				"client_ip": c.ClientIP(),
				"path":      c.Request.URL.Path,
			}).Warn("Akses ditolak: Token di-blacklist atau sudah di-logout")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak valid atau sudah di-logout",
			})
			c.Abort()
			return
		}

		// Validasi token
		_, claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.Log.WithFields(map[string]interface{}{
				"client_ip": c.ClientIP(),
				"path":      c.Request.URL.Path,
				"error":     err.Error(),
			}).Warn("Akses ditolak: Validasi Token JWT gagal")
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token tidak valid: " + err.Error(),
			})
			c.Abort()
			return
		}

		// Simpan claims di context
		if userID, ok := claims["user_id"]; ok {
			c.Set("userID", userID)
		}
		if email, ok := claims["email"]; ok {
			c.Set("userEmail", email)
		}
		if role, ok := claims["role"]; ok {
			c.Set("userRole", role)
		}

		c.Next()
	}
}
