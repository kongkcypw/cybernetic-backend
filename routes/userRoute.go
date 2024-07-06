package routes

import (
	controllers "example/backend/controllers"
	"example/backend/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.GET("user/character", middleware.Authenticate(), controllers.GetCharacter())
	router.GET("user/game-play-level", controllers.GetUserGamePlayLevel())
	router.POST("user/game-play-level/update", controllers.UpdateLevelPlayed())
}
