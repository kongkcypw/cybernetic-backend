package routes

import (
	controllers "example/backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/user/signup", controllers.Signup())
	router.GET("/user/info", controllers.GetUser())
}
