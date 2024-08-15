package routes

import (
	controllers "example/backend/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/user/character/update", controllers.UpdateCharacter())
	router.GET("/user/character", controllers.GetCharacter())
	router.POST("/user/character/update/level", controllers.UpdateHighestLevel())

	// Not used
	router.GET("/user/game-play-level", controllers.GetUserGamePlayLevel())
	router.POST("/user/game-play-level/update", controllers.UpdateLevelPlayed())
}
