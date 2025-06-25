// @title Inventory Management API
// @version 1.0
// @description RESTful API untuk mengelola inventaris (kategori, produk, stok, dan user).
// @termsOfService http://example.com/terms/

// @contact.name Danar Rafiardi
// @contact.email danarbaadilla12@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"inventory-management-api/app"
	"inventory-management-api/config"
	"inventory-management-api/controller"
	"inventory-management-api/repository"
	"inventory-management-api/route"
	"inventory-management-api/service"

	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	// Tambahan untuk Swagger
	_ "inventory-management-api/docs"

	swagger "github.com/swaggo/fiber-swagger"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("[WARNING] .env file not found, using default env")
	}

	// Inisialisasi koneksi database
	db, err := config.NewGormMySQLConnection()
	if err != nil {
		log.Fatalf("❌ Gagal konek database: %v", err)
	}

	// Inisialisasi validator
	validate := validator.New()

	// Inisialisasi repository
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	stockMovementRepo := repository.NewStockMovementRepository(db)

	// Inisialisasi service
	authService := service.NewAuthService(userRepo)
	userService := service.NewUserService(userRepo, validate)
	categoryService := service.NewCategoryService(categoryRepo, validate)
	productService := service.NewProductService(productRepo, validate)
	stockMovementService := service.NewStockMovementService(stockMovementRepo, productRepo, db, validate)

	// Inisialisasi controller
	authController := controller.NewAuthController(authService, userService)
	userController := controller.NewUserController(userService)
	categoryController := controller.NewCategoryController(categoryService)
	productController := controller.NewProductController(productService)
	stockMovementController := controller.NewStockMovementController(stockMovementService)

	// Inisialisasi Fiber app
	fiberApp := app.NewApp()

	// Aktifkan CORS untuk semua origin (bisa dibatasi jika sudah production)
	fiberApp.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Development mode, menerima semua akses (untuk Portfolio)
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Endpoint Swagger UI
	fiberApp.Get("/swagger/*", swagger.FiberWrapHandler())

	// Registrasi semua routes
	route.RegisterAuthRoutes(fiberApp, authController)
	route.RegisterUserRoutes(fiberApp, userController)
	route.RegisterCategoryRoutes(fiberApp, categoryController)
	route.RegisterProductRoutes(fiberApp, productController)
	route.RegisterStockMovementRoutes(fiberApp, stockMovementController)

	// Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Printf("🚀 Server berjalan di http://localhost:%s", port)
	log.Printf("📚 Dokumentasi Swagger tersedia di http://localhost:%s/swagger/index.html", port)
	log.Fatal(fiberApp.Listen(":" + port))
}
