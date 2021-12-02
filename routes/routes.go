package routes

import (
	"concierge/config"
	"concierge/controllers"
	"concierge/routes/middleware"

	"github.com/gin-gonic/gin"
)

//InitializeRoutes ...
func InitializeRoutes(router *gin.Engine) {

	// Handle the index route
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "http://"+c.Request.Host+"/resources")
	})

	authorized := router.Group("/")
	authorized.Use(middleware.Authorize)
	{
		authorized.GET("/resources", controllers.ShowAllowedIngress)
		authorized.GET("/resources/:driver/:ns/:name", controllers.IngressDetails)
		authorized.POST("/resources/:driver/:ns/:name", controllers.WhiteListIP)
		authorized.POST("/resources/:driver/:ns/:name/:id", controllers.DeleteIPFromIngress)
		authorized.GET("/users", controllers.GetUsers)
		authorized.GET("/users/add", controllers.AddUsersForm)
		authorized.POST("/users/add", controllers.AddUsers)
		authorized.GET("/users/edit/:id", controllers.UpdateUsersForm)
		authorized.POST("/users/edit/:id", controllers.UpdateUser)
		authorized.POST("/users/delete", controllers.DeleteUser)
		authorized.GET("/logout", func(c *gin.Context) {
			c.Redirect(302, "http://"+c.Request.Host+"/oauth2/sign_in")
		})
	}

	cron := router.Group("/cron", gin.BasicAuth(gin.Accounts{
		config.CronConfig.CronUsername: config.CronConfig.CronPassword,
	}))
	cron.Use(middleware.Cron)
	{
		cron.GET("/", controllers.ClearExpiredLeases)
		cron.POST("/", controllers.ClearExpiredLeases)
	}
}
