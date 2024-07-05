package routes

import (
	controllers "example/backend/controllers"

	"github.com/gin-gonic/gin"
)

func FirebaseRoutes(router *gin.Engine) {
	router.POST("/firebase/upload-image", controllers.UploadImageToFirebase())
}
