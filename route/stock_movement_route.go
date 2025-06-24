package route

import (
	"inventory-management-api/controller"
	"inventory-management-api/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterStockMovementRoutes(app *fiber.App, c *controller.StockMovementController) {
	// Group untuk stock movements
	stock := app.Group("/stock-movements", middleware.JWTMiddleware)

	// Dapat diakses oleh admin dan staff
	stock.Get("/", c.FindAll)
	stock.Get("/:id", c.FindById)

	// Hanya Staff yang boleh buat transaksi
	stock.Post("/", middleware.StaffOnly, c.Create)

	// Hanya Admin yang boleh hapus transaksi
	stock.Delete("/:id", middleware.AdminOnly, c.Delete)

	// ✅ Endpoint laporan bulanan - hanya admin yang boleh akses
	app.Get("/reports/stock-movements", middleware.JWTMiddleware, middleware.AdminOnly, c.GetMonthlyReport)
}
