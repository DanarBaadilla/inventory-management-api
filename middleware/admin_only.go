package middleware

import (
	"inventory-management-api/model/web"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func AdminOnly(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "admin" {
		return c.Status(http.StatusForbidden).JSON(web.WebResponse{
			Code:   http.StatusForbidden,
			Status: "FORBIDDEN",
			Error:  "Admin only access",
		})
	}
	return c.Next()
}
