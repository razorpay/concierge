package controllers

import (
	"concierge/database"
	"concierge/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

//GetUsers ...
func GetUsers(c *gin.Context) {
	User, _ := c.Get("User")
	db, err := database.Conn()
	if err != nil {
		log.Error("Error", err)
	}
	defer db.Close()

	users := []models.Users{}
	db.Find(&users)
	c.HTML(http.StatusOK, "manageusers.gohtml", gin.H{
		"data": users,
		"user": User,
	})
}

//AddUsersForm ...
func AddUsersForm(c *gin.Context) {
	User, _ := c.Get("User")
	c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
		"user": User,
	})
}

//AddUsers ...
func AddUsers(c *gin.Context) {
	User, _ := c.Get("User")
	var newUser models.Users
	if err := c.ShouldBind(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := database.Conn()
	if err != nil {
		log.Error("Error", err)
	}
	defer db.Close()
	split := strings.Split(newUser.Email, "@")
	username := split[0]
	if split[1] != "razorpay.com" {
		c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": "Invalid Organization Email",
			},
		})
		return
	}
	myUser := models.Users{
		Username: username,
		Email:    newUser.Email,
	}
	newUser.Username = username
	res := db.Where(myUser).First(&models.Users{})
	if res.RecordNotFound() {
		db.Create(&newUser)
		c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
			"user": User,
			"message": map[string]string{
				"class":   "Success",
				"message": "User is successfully created",
			},
		})
		return
	}
	c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
		"user": User,
		"message": map[string]string{
			"class":   "Warning",
			"message": "User is already present",
		},
	})
}
