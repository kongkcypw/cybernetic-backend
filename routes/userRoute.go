package routes

import (
	controllers "example/backend/controllers"
	"example/backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(router *fiber.App) {
	router.Get("/user/info", middleware.Authenticate, controllers.GetUser)
	router.Get("/user/character", middleware.Authenticate, controllers.GetCharacter)
}
