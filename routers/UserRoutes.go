package routers

import (
	"RemiAPI/controllers"
	"RemiAPI/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes sets up authentication routes
func UserRoutes(router *gin.Engine) {
	userGroup := router.Group("/auth")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.POST("/me", controllers.GetUserDetailsHandler)
		userGroup.POST("/update-profile", controllers.UpdateProfileHandler)
	}
}
