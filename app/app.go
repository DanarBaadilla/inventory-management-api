package app

import (
	"inventory-management-api/model/web"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			statusCode := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := err.(*fiber.Error); ok {
				statusCode = e.Code
				message = e.Message
			} else if err != nil {
				log.Printf("[ERROR] Unhandled error: %v\n", err)
			}

			return c.Status(statusCode).JSON(web.WebResponse{
				Code:   statusCode,
				Status: http.StatusText(statusCode),
				Data:   nil,
				Error:  message,
			})
		},
	})

	return app
}
