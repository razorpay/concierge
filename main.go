package main

import (
	"concierge/config"
	"concierge/database"
	"concierge/models"
	"concierge/routes"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var router *gin.Engine

func main() {
	config.LoadConfig()
	database.Conn()
	defer database.CloseDB()
	migrations()
	seeding()
	// Start the router
	router = gin.Default()
	// Serving static files
	router.Static("/assets", "./assets")

	router.LoadHTMLGlob("templates/*")

	// Initialize the routes
	routes.InitializeRoutes(router)
	listenPort := os.Getenv("APP_PORT")
	router.Run("0.0.0.0:" + listenPort)

}

func migrations() {
	if database.DB == nil {
		database.Conn()
	}
	database.DB.AutoMigrate(&models.Users{}, &models.Leases{})
	database.DB.Model(&models.Leases{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
}

func seeding() {
	database.Seeding()
}
