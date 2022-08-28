package middleware

import (
	"concierge/config"
	"encoding/base64"
	"strings"

	"github.com/gin-gonic/gin"
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
	username := pair[0]
	password := pair[1]
	if username != config.CronConfig.CronUsername {
		c.AbortWithStatusJSON(401, "Invalid username.")
		return
	}
	if password != config.CronConfig.CronPassword {
		c.AbortWithStatusJSON(401, "Invalid password.")
		return
	}
	c.Next()
}
