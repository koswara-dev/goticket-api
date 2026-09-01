package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gotiket-api/middleware"

	"github.com/gin-gonic/gin"
)

func TestRequiredRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		allowedRoles   []string
		contextRole    interface{}
		setRole        bool
		expectedStatus int
	}{
		{
			name:           "Success - Single Allowed Role",
			allowedRoles:   []string{"admin"},
			contextRole:    "admin",
			setRole:        true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Success - Multiple Allowed Roles",
			allowedRoles:   []string{"admin", "superadmin"},
			contextRole:    "superadmin",
			setRole:        true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Forbidden - Role Not Allowed",
			allowedRoles:   []string{"admin"},
			contextRole:    "customer",
			setRole:        true,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "Unauthorized - Role Not In Context",
			allowedRoles:   []string{"admin"},
			contextRole:    nil,
			setRole:        false,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Unauthorized - Invalid Role Type",
			allowedRoles:   []string{"admin"},
			contextRole:    12345, // int instead of string
			setRole:        true,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()

			// Middleware dummy untuk menyuntikkan userRole ke context
			router.Use(func(c *gin.Context) {
				if tt.setRole {
					c.Set("userRole", tt.contextRole)
				}
				c.Next()
			})

			// Pasang RequiredRole middleware yang diuji
			router.GET("/protected", middleware.RequiredRole(tt.allowedRoles...), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest(http.MethodGet, "/protected", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}
