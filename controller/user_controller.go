package controller

import (
	"inventory-management-api/model/web"
	"inventory-management-api/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{UserService: userService}
}

// FindAll godoc
// @Summary Ambil semua data user
// @Description Endpoint ini digunakan untuk mengambil semua user yang terdaftar dalam sistem.
// @Tags User
// @Produce json
// @Success 200 {object} web.WebResponse{data=[]web.UserResponse}
// @Failure 500 {object} web.WebResponse
// @Security BearerAuth
// @Router /users [get]
func (c *UserController) FindAll(ctx *fiber.Ctx) error {
	users, err := c.UserService.FindAll()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Error:  err.Error(),
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   users,
	})
}

// FindByID godoc
// @Summary Ambil user berdasarkan ID
// @Description Endpoint ini digunakan untuk mendapatkan data user berdasarkan ID.
// @Tags User
// @Produce json
// @Param id path int true "ID user yang ingin diambil"
// @Success 200 {object} web.WebResponse{data=web.UserResponse}
// @Failure 400,404 {object} web.WebResponse
// @Security BearerAuth
// @Router /users/{id} [get]
func (c *UserController) FindByID(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid user ID",
		})
	}

	user, err := c.UserService.FindByID(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Error:  "User not found",
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   user,
	})
}

// Create godoc
// @Summary Tambah user baru
// @Description Endpoint ini digunakan untuk menambahkan user baru ke sistem.
// @Tags User
// @Accept json
// @Produce json
// @Param request body web.UserCreateOrUpdateRequest true "Data user baru"
// @Success 201 {object} web.WebResponse{data=web.UserResponse}
// @Failure 400,500 {object} web.WebResponse
// @Security BearerAuth
// @Router /users [post]
func (c *UserController) Create(ctx *fiber.Ctx) error {
	var req web.UserCreateOrUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid request body",
		})
	}

	user, err := c.UserService.Create(req)
	if err != nil {
		if strings.HasPrefix(err.Error(), "validation error:") {
			return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "BAD REQUEST",
				Error:  err.Error(),
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Error:  err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(web.WebResponse{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   user,
	})
}

// Update godoc
// @Summary Perbarui data user
// @Description Endpoint ini digunakan untuk memperbarui informasi user berdasarkan ID.
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "ID user yang akan diperbarui"
// @Param request body web.UserCreateOrUpdateRequest true "Data user terbaru"
// @Success 200 {object} web.WebResponse{data=web.UserResponse}
// @Failure 400,404,500 {object} web.WebResponse
// @Security BearerAuth
// @Router /users/{id} [put]
func (c *UserController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid user ID",
		})
	}

	var req web.UserCreateOrUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid request body",
		})
	}

	user, err := c.UserService.Update(id, req)
	if err != nil {
		if err.Error() == "user not found" {
			return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
				Code:   http.StatusNotFound,
				Status: "NOT FOUND",
				Error:  "User not found",
			})
		}
		if strings.HasPrefix(err.Error(), "validation error:") {
			return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "BAD REQUEST",
				Error:  err.Error(),
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Error:  err.Error(),
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   user,
	})
}

// Delete godoc
// @Summary Hapus user berdasarkan ID
// @Description Endpoint ini digunakan untuk menghapus user dari sistem berdasarkan ID-nya.
// @Tags User
// @Produce json
// @Param id path int true "ID user yang ingin dihapus"
// @Success 200 {object} web.WebResponse{data=string}
// @Failure 400,404 {object} web.WebResponse
// @Security BearerAuth
// @Router /users/{id} [delete]
func (c *UserController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid user ID",
		})
	}

	err = c.UserService.Delete(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Error:  "User not found",
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "User deleted",
	})
}
