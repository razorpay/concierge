package middleware

import (
	"github.com/gin-gonic/gin"
)

// SetCustomHeaders ...
func SetCustomHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Next()
	}
}
