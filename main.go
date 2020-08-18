package main

import (
	"concierge/config"
	"concierge/database"
	"concierge/models"
	"concierge/routes"
	"concierge/routes/middleware"

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
	// Creates a router without any middleware by default
	router := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(middleware.Recovery)
	router.Use(middleware.CSRF())
	router.Use(middleware.SetCustomHeaders())

	// Serving static files
	router.Static("/assets", "./assets")

	router.LoadHTMLGlob("templates/*")

	// Initialize the routes
	routes.InitializeRoutes(router)
	router.Run(config.AppCfg.ListenIP + ":" + config.AppCfg.ListenPort)

}

func migrations() {
	if database.DB == nil {
		database.Conn()
	}

	database.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci").AutoMigrate(&models.Users{}, &models.Leases{})
	database.DB.Model(&models.Leases{}).AddForeignKey("user_id", "users(id)", "CASCADE", "RESTRICT")
}

func seeding() {
	if config.AppCfg.Mode == "dev" {
		database.Seeding()
	}
}
