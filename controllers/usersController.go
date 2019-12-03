package controllers

import (
	"concierge/database"
	"concierge/models"
	"net/http"
	"os"
	"strconv"
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

//UpdateUsersForm ...
func UpdateUsersForm(c *gin.Context) {
	User, _ := c.Get("User")
	ID := c.Param("id")
	db, err := database.Conn()
	if err != nil {
		log.Error("Error", err)
	}
	defer db.Close()
	editUser := models.Users{}
	db.Find(&editUser, ID)
	c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
		"user": User,
		"data": editUser,
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
	if split[1] != os.Getenv("COMPANY_DOMAIN") {
		c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": "Invalid Organization Email",
			},
			"data": newUser,
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
			"data": newUser,
		})
		return
	}
	c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
		"user": User,
		"message": map[string]string{
			"class":   "Warning",
			"message": "User is already present",
		},
		"data": newUser,
	})
}

//UpdateUser ...
func UpdateUser(c *gin.Context) {
	User, _ := c.Get("User")
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	ID := uint(id)
	var updateUser models.Users
	if err := c.ShouldBind(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateUser.ID = ID
	db, err := database.Conn()
	if err != nil {
		log.Error("Error", err)
	}
	defer db.Close()
	split := strings.Split(updateUser.Email, "@")
	username := split[0]
	if split[1] != os.Getenv("COMPANY_DOMAIN") {
		c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": "Invalid Organization Email",
			},
			"data": updateUser,
		})
		return
	}
	myUser := models.Users{
		Username: username,
		Email:    updateUser.Email,
	}
	updateUser.Username = username
	res := db.Where(myUser).Where("id != ?", ID).First(&models.Users{})
	if res.RecordNotFound() {
		log.Info(updateUser)
		db.Model(&updateUser).Updates(updateUser)
		c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
			"user": User,
			"message": map[string]string{
				"class":   "Success",
				"message": "User is updated created",
			},
			"data": updateUser,
		})
		return
	}
	c.HTML(http.StatusOK, "addusers.gohtml", gin.H{
		"user": User,
		"message": map[string]string{
			"class":   "Warning",
			"message": "User is already present",
		},
		"data": updateUser,
	})
}

//DeleteUser ...
func DeleteUser(c *gin.Context) {
	User, _ := c.Get("User")
	id, _ := strconv.ParseUint(c.PostForm("ID"), 10, 64)
	ID := uint(id)
	db, err := database.Conn()
	if err != nil {
		log.Error("Error", err)
	}
	defer db.Close()

	users := []models.Users{}
	db.Find(&users)
	if User.(*models.Users).ID == ID {
		c.HTML(http.StatusOK, "manageusers.gohtml", gin.H{
			"user": User,
			"message": map[string]string{
				"class":   "Warning",
				"message": "You can't delete yourself",
			},
			"data": users,
		})
		return
	}

	res := db.First(&models.Users{}, ID)
	if res.RecordNotFound() {
		c.HTML(http.StatusOK, "manageusers.gohtml", gin.H{
			"user": User,
			"message": map[string]string{
				"class":   "Danger",
				"message": "User is not found",
			},
			"data": users,
		})
		return
	}
	db.Delete(models.Users{}, "id = ?", ID)
	db.Find(&users)
	c.HTML(http.StatusOK, "manageusers.gohtml", gin.H{
		"user": User,
		"message": map[string]string{
			"class":   "Success",
			"message": "User is deleted successfully",
		},
		"data": users,
	})
}
