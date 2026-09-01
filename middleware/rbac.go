package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// middleware otorisasi peran, return response success, message, data
func RequiredRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// dapatkan data role dari context gin (suntikan dari auth middleware)
		roleVal, exist := c.Get("userRole")

		if !exist {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Akses ditolak. Anda tidak terautentikasi.",
			})
			return
		}

		role, ok := roleVal.(string)
		if !ok {
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
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Akses ditolak. Anda tidak memiliki hak akses.",
			})
			return
		}

		c.Next()
	}
}
