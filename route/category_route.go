package route

import (
	"inventory-management-api/controller"
	"inventory-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterCategoryRoutes(app *fiber.App, controller *controller.CategoryController) {
	category := app.Group("/categories", middleware.JWTMiddleware)

	// Bisa diakses oleh admin dan staff
	category.Get("/", controller.FindAll)
	category.Get("/:id", controller.FindById)

	// Hanya admin yang boleh manipulasi data kategori
	category.Post("/", middleware.AdminOnly, controller.Create)
	category.Put("/:id", middleware.AdminOnly, controller.Update)
	category.Delete("/:id", middleware.AdminOnly, controller.Delete)
}
