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
		authorized.GET("/logout", func(c *gin.Context) {
			c.SetCookie("_oauth2_proxy", "", -1, "/", "https://"+c.Request.Host+"/", false, true)
			c.Redirect(302, "http://"+c.Request.Host+"/")
		})
	}
}
