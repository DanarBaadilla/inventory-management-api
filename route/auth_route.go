package route

import (
	"inventory-management-api/controller"
	"inventory-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app *fiber.App, controller *controller.AuthController) {
	// Endpoint login (tanpa middleware)
	app.Post("/login", controller.Login)

	// Group untuk endpoint yang butuh JWT
	auth := app.Group("/auth", middleware.JWTMiddleware)
	auth.Get("/me", controller.Me)
}
