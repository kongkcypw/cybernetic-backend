package routes

import (
	controllers "example/backend/controllers"
	"example/backend/middleware"

	"github.com/gin-gonic/gin"
)

func CharacterRoutes(router *gin.Engine) {
	router.GET("/character/info", middleware.Authenticate(), controllers.GetCharacter())
}
