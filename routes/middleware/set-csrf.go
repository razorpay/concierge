package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

// SetXCSRFToken ...
func SetXCSRFToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-CSRF-Token", csrf.Token(c.Request))
		c.Next()
	}
}
