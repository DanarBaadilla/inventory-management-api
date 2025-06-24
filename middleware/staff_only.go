package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func StaffOnly(c *fiber.Ctx) error {
	role, ok := c.Locals("role").(string)
	if !ok || role != "staff" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":   fiber.StatusForbidden,
			"status": "FORBIDDEN",
			"error":  "Staff access only",
		})
	}
	return c.Next()
}
