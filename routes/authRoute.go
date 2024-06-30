package routes

import (
	controllers "example/backend/controllers"
	middleware "example/backend/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/auth/signup", controllers.Signup())
	router.POST("/auth/login", controllers.Login())
	router.GET("/auth/logout", controllers.Logout())
	// Google OAuth
	router.GET("/auth/google/login", middleware.GoogleCallback(), controllers.LoginWithGoogle())
}
