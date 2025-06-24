package route

import (
	"inventory-management-api/controller"
	"inventory-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(app *fiber.App, controller *controller.ProductController) {
	product := app.Group("/products", middleware.JWTMiddleware)

	// Boleh diakses oleh staff dan admin
	product.Get("/search", controller.Search)
	product.Get("/", controller.FindAll)
	product.Get("/:id", controller.FindById)

	// Hanya admin yang boleh manipulasi data
	product.Post("/", middleware.AdminOnly, controller.Create)
	product.Put("/:id", middleware.AdminOnly, controller.Update)
	product.Delete("/:id", middleware.AdminOnly, controller.Delete)
}
