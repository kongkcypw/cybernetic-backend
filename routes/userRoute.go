package routes

import (
	controllers "example/backend/controllers"
	"example/backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("/user/info", middleware.Authenticate(), controllers.GetUser())
}
