package controller

import (
	"inventory-management-api/model/web"
	"inventory-management-api/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	Service service.CategoryService
}

func NewCategoryController(service service.CategoryService) *CategoryController {
	return &CategoryController{Service: service}
}

// FindAll godoc
// @Summary Mendapatkan semua kategori
// @Description Mengambil semua data kategori yang tersedia
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Success 200 {object} web.WebResponse{data=[]web.CategoryResponse}
// @Failure 500 {object} web.WebResponse
// @Router /categories [get]
func (c *CategoryController) FindAll(ctx *fiber.Ctx) error {
	result, err := c.Service.FindAll()
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
		Data:   result,
	})
}

// FindById godoc
// @Summary Mendapatkan kategori berdasarkan ID
// @Description Mengambil detail kategori berdasarkan ID yang diberikan
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Kategori"
// @Success 200 {object} web.WebResponse{data=web.CategoryResponse}
// @Failure 400 {object} web.WebResponse
// @Failure 404 {object} web.WebResponse
// @Router /categories/{id} [get]
func (c *CategoryController) FindById(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid category ID",
		})
	}

	result, err := c.Service.FindById(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Error:  "Category not found",
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	})
}

// Create godoc
// @Summary Membuat kategori baru
// @Description Menambahkan kategori baru ke dalam sistem
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body web.CategoryCreateOrUpdateRequest true "Data kategori baru"
// @Success 201 {object} web.WebResponse{data=web.CategoryResponse}
// @Failure 400 {object} web.WebResponse
// @Router /categories [post]
func (c *CategoryController) Create(ctx *fiber.Ctx) error {
	var req web.CategoryCreateOrUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid request body",
		})
	}

	result, err := c.Service.Create(req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  err.Error(),
		})
	}

	return ctx.Status(http.StatusCreated).JSON(web.WebResponse{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   result,
	})
}

// Update godoc
// @Summary Memperbarui kategori
// @Description Mengubah data kategori berdasarkan ID
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Kategori"
// @Param request body web.CategoryCreateOrUpdateRequest true "Data kategori yang diperbarui"
// @Success 200 {object} web.WebResponse{data=web.CategoryResponse}
// @Failure 400 {object} web.WebResponse
// @Router /categories/{id} [put]
func (c *CategoryController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid category ID",
		})
	}

	var req web.CategoryCreateOrUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid request body",
		})
	}

	result, err := c.Service.Update(id, req)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  err.Error(),
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	})
}

// Delete godoc
// @Summary Menghapus kategori
// @Description Menghapus kategori berdasarkan ID
// @Tags Categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID Kategori"
// @Success 200 {object} web.WebResponse{data=string}
// @Failure 400 {object} web.WebResponse
// @Failure 404 {object} web.WebResponse
// @Router /categories/{id} [delete]
func (c *CategoryController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid category ID",
		})
	}

	err = c.Service.Delete(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Error:  "Category not found",
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Category deleted",
	})
}
