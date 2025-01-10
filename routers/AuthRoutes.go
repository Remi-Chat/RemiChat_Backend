package routers

import (
	"RemiAPI/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up authentication routes
func AuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", controllers.SignupHandler)
		authGroup.POST("/login", controllers.LoginHandler)
		authGroup.POST("/temp-user", controllers.CreateTempUser)
	}
}
