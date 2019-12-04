package middleware

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//Cron ...
func Cron(c *gin.Context) {
	auth := strings.SplitN(c.Request.Header.Get("Authorization"), " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		c.AbortWithStatusJSON(401, "Unauthorized")
		return
	}
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	log.Info(c.Request.Header.Get("Authorization"))

	username := pair[0]
	password := pair[1]
	if username != os.Getenv("CRON_USERNAME") {
		c.AbortWithStatusJSON(401, "Invalid username.")
		return
	}
	if password != os.Getenv("CRON_PASSWORD") {
		c.AbortWithStatusJSON(401, "Invalid password.")
		return
	}
	c.Next()
}
