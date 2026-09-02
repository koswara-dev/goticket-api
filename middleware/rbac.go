package middleware

import (
	"net/http"

	"gotiket-api/utils"

	"github.com/gin-gonic/gin"
)

// middleware otorisasi peran, return response success, message, data
func RequiredRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// dapatkan data role dari context gin (suntikan dari auth middleware)
		roleVal, exist := c.Get("userRole")

		if !exist {
			utils.Log.WithFields(map[string]interface{}{
				"client_ip": c.ClientIP(),
				"path":      c.Request.URL.Path,
			}).Warn("Akses RBAC ditolak: Data role tidak ditemukan di context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Akses ditolak. Anda tidak terautentikasi.",
			})
			return
		}

		role, ok := roleVal.(string)
		if !ok {
			utils.Log.WithFields(map[string]interface{}{
				"client_ip": c.ClientIP(),
				"path":      c.Request.URL.Path,
			}).Warn("Akses RBAC ditolak: Format role di context tidak valid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Akses ditolak. Format role tidak valid.",
			})
			return
		}

		isAllowed := false

		for _, r := range roles {
			if role == r {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			utils.Log.WithFields(map[string]interface{}{
				"client_ip":     c.ClientIP(),
				"path":          c.Request.URL.Path,
				"user_role":     role,
				"allowed_roles": roles,
			}).Warn("Akses RBAC ditolak: Hak akses (Role) tidak mencukupi")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Anda tidak memiliki hak akses.",
			})
			return
		}

		c.Next()
	}
}
