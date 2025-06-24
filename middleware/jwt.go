package middleware

import (
	"inventory-management-api/helper"
	"inventory-management-api/model/web"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware(c *fiber.Ctx) error {
	// Ambil header Authorization
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Error:  "Missing Authorization header",
		})
	}

	// Format harus "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(http.StatusUnauthorized).JSON(web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Error:  "Invalid Authorization header format",
		})
	}

	tokenString := parts[1]

	// Validasi token dan ambil claims
	claims, err := helper.ValidateToken(tokenString)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Error:  "Invalid or expired token",
		})
	}

	// Set user_id dan role ke context
	c.Locals("user_id", claims.UserID)
	c.Locals("role", claims.Role)

	return c.Next()
}
