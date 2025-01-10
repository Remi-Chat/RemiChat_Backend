package routers

import (
	"RemiAPI/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up authentication routes
func RegisterAuthRoutes(router *gin.Engine) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/signup", controllers.SignupHandler)
		authGroup.POST("/login", controllers.LoginHandler)
	}
}
