package middleware

import (
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Recovery : Middlewares to catch recovers from any panics
// and writes a 500 if there was one.
func Recovery(c *gin.Context) {
	defer func(c *gin.Context) {
		handlePanic(c)
	}(c)

	c.Next()
}

func handlePanic(c *gin.Context) {
	if err := recover(); err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe bool
		if ne, ok := err.(*net.OpError); ok {
			if se, ok := ne.Err.(*os.SyscallError); ok {
				if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
					brokenPipe = true
				}
			}
		}

		if brokenPipe {
			log.Errorf("%s", err)
			c.Error(err.(error)) // nolint: errcheck
			c.Abort()
		} else {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}
}
