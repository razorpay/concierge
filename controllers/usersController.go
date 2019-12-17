package controllers

import (
	"concierge/database"
	"concierge/models"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//GetUsers ...
func GetUsers(c *gin.Context) {
	User, _ := c.Get("User")
	if database.DB == nil {
		database.Conn()
	}

	users := []models.Users{}
	database.DB.Find(&users)
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
	if database.DB == nil {
		database.Conn()
	}
	editUser := models.Users{}
	database.DB.Find(&editUser, ID)
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
	if database.DB == nil {
		database.Conn()
	}
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
	newUser.Name = username
	res := database.DB.Where(myUser).First(&models.Users{})
	if res.RecordNotFound() {
		database.DB.Create(&newUser)
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
	if database.DB == nil {
		database.Conn()
	}
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
	updateUser.Name = username
	res := database.DB.Where(myUser).Where("id != ?", ID).First(&models.Users{})
	if res.RecordNotFound() {
		database.DB.Model(&updateUser).Updates(updateUser)
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
	if database.DB == nil {
		database.Conn()
	}

	users := []models.Users{}
	database.DB.Find(&users)
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

	res := database.DB.First(&models.Users{}, ID)
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
	database.DB.Delete(models.Users{}, "id = ?", ID)
	database.DB.Find(&users)
	c.HTML(http.StatusOK, "manageusers.gohtml", gin.H{
		"user": User,
		"message": map[string]string{
			"class":   "Success",
			"message": "User is deleted successfully",
		},
		"data": users,
	})
}
