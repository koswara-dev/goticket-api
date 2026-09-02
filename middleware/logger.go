package middleware

import (
	"time"

	"gotiket-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GinLogger middleware mencatat detail lalu lintas HTTP request ke utils.Log secara terstruktur.
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.URL.Path
		rawQuery := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(startTime)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if rawQuery != "" {
			path = path + "?" + rawQuery
		}

		fields := logrus.Fields{
			"status":     statusCode,
			"latency_ms": float64(latency.Nanoseconds()) / 1e6,
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"user_agent": c.Request.UserAgent(),
		}

		if errorMessage != "" {
			fields["error"] = errorMessage
		}

		entry := utils.Log.WithFields(fields)

		if statusCode >= 500 {
			entry.Error("HTTP Server Error")
		} else if statusCode >= 400 {
			entry.Warn("HTTP Client Warning")
		} else {
			entry.Info("HTTP Request Handled")
		}
	}
}
