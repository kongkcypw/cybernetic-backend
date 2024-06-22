package routes

import (
	controllers "example/backend/controllers"
	middleware "example/backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(router *fiber.App) {
	router.Post("/auth/signup", controllers.Signup)
	router.Post("/auth/login", controllers.Login)
	router.Get("/auth/logout", controllers.Logout)
	// Google OAuth
	router.Get("/auth/google/signup", middleware.GoogleCallback, controllers.SignupWithGoogle)
	router.Get("/auth/google/login", middleware.GoogleCallback, controllers.LoginWithGoogle)
}
