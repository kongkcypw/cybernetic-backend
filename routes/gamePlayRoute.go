package routes

import (
	controllers "example/backend/controllers"

	"github.com/gin-gonic/gin"
)

func GamePlayRoutes(router *gin.Engine) {
	router.GET("/game-play/level-selection", controllers.GetGamePlayLevel())
	router.GET("/game-play/level-selection/get-one", controllers.GetGamePlayLevelByNumber())
}
