package routes

import (
	"concierge/controllers"
	"concierge/routes/middleware"

	"github.com/gin-gonic/gin"
)

//InitializeRoutes ...
func InitializeRoutes(router *gin.Engine) {

	// Handle the index route
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "http://"+c.Request.Host+"/ingress")
	})

	authorized := router.Group("/")
	authorized.Use(middleware.Authorize)
	{
		authorized.GET("/ingress", controllers.ShowAllowedIngress)
		authorized.GET("/ingress/:ns/:name", controllers.IngressDetails)
		authorized.POST("/ingress/:ns/:name", controllers.WhiteListIP)
		authorized.POST("/ingress/:ns/:name/:id", controllers.DeleteIPFromIngress)
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

	router.GET("/cron", controllers.ClearExpiredLeases)
	cron := router.Group("/cron")
	cron.Use(middleware.Cron)
	{
		cron.POST("/", controllers.ClearExpiredLeases)
	}
}
