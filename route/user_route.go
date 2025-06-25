package route

import (
	"inventory-management-api/controller"
	"inventory-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App, controller *controller.UserController) {
	userGroup := app.Group("/users", middleware.JWTMiddleware, middleware.AdminOnly)

	userGroup.Get("/", controller.FindAll)
	userGroup.Get("/:id", controller.FindByID)
	userGroup.Post("/", controller.Create)
	userGroup.Put("/:id", controller.Update)
	userGroup.Delete("/:id", controller.Delete)
}
