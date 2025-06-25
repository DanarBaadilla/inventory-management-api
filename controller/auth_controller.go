package controller

import (
	"inventory-management-api/model/web"
	"inventory-management-api/service"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	AuthService service.AuthService
	UserService service.UserService
}

func NewAuthController(authService service.AuthService, userService service.UserService) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserService: userService,
	}
}

// Login godoc
// @Summary Login untuk mendapatkan token JWT
// @Description Autentikasi user berdasarkan email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body web.LoginRequest true "Login credentials"
// @Success 200 {object} web.WebResponse{data=web.LoginResponse}
// @Failure 400 {object} web.WebResponse
// @Failure 401 {object} web.WebResponse
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req web.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid request body",
		})
	}

	resp, err := c.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Error:  err.Error(),
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   resp,
	})
}

// Me godoc
// @Summary Mendapatkan informasi user yang sedang login
// @Description Endpoint ini membutuhkan token JWT yang valid
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} web.WebResponse{data=web.UserResponse}
// @Failure 404 {object} web.WebResponse
// @Router /auth/me [get]
func (c *AuthController) Me(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(int)

	user, err := c.UserService.FindByID(userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(web.WebResponse{
			Code:   fiber.StatusNotFound,
			Status: "NOT FOUND",
			Error:  "User not found",
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   200,
		Status: "OK",
		Data:   user,
	})
}
