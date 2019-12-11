package middleware

import (
	"concierge/database"
	"concierge/models"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//Authorize ...
func Authorize(c *gin.Context) {
	// Disable user groups and enable user email
	db, err := database.Conn()
	if err != nil {
		log.Error("Error", err)
	}
	defer db.Close()
	username := c.GetHeader("X-Forwarded-User")
	email := c.GetHeader("X-Forwarded-Email")
	split := strings.Split(email, "@")
	if split[1] != os.Getenv("COMPANY_DOMAIN") {
		c.AbortWithStatusJSON(404, "Invalid Organization Email")
		return
	}
	user := models.Users{
		Username: username,
		Name:     username,
		Email:    email,
	}
	getUser := &models.Users{}
	db.FirstOrCreate(getUser, user)
	c.Set("User", getUser)
	c.Next()
}
