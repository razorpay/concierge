package main

import (
	"concierge/config"
	"concierge/database"
	"concierge/models"
	"concierge/routes"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

var router *gin.Engine

func main() {
	config.LoadConfig()
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
	router.Run(":" + listenPort)

}

func migrations() {
	var db *gorm.DB
	db, err := database.Conn()
	if err != nil {
		log.Error(err)
	}
	db.AutoMigrate(&models.Users{}, &models.Leases{})
	db.Model(&models.Leases{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
	defer db.Close()
}

func seeding() {
	var db *gorm.DB
	db, err := database.Conn()
	if err != nil {
		log.Error(err)
	}
	db.FirstOrCreate(&models.Users{}, &models.Users{
		Name:     "Ankit Jain",
		Username: "ankit.infra",
		Email:    "ankit.infra@razorpay.com",
		Admin:    1,
	})
	defer db.Close()
}
