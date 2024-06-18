package routes

import (
	controllers "example/backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router *fiber.App) {
	router.Post("/auth/signup", controllers.Signup)
	router.Post("/auth/login", controllers.Login)
	router.Get("/auth/logout", controllers.Logout)
	router.Get("/auth/google/login", controllers.GoogleLogin)
	router.Get("/auth/google/callback", controllers.GoogleCallback)
}
