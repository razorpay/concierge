package middleware

import (
	"concierge/database"
	"concierge/models"

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
	user := models.Users{
		Username: username,
		Email:    email,
	}
	getUser := &models.Users{}
	res := db.Where(user).First(getUser)
	if res.RecordNotFound() {
		log.Info("Record not found")
		c.AbortWithStatusJSON(404, "Not Found")
		return
	}
	c.Set("User", getUser)
	c.Next()
}
