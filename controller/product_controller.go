package controller

import (
	"inventory-management-api/model/web"
	"inventory-management-api/service"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	Service service.ProductService
}

func NewProductController(service service.ProductService) *ProductController {
	return &ProductController{Service: service}
}

// FindAll godoc
// @Summary Mendapatkan Seluruh Produk
// @Description Mengambil seluruh data produk pada database
// @Tags Product
// @Produce json
// @Security BearerAuth
// @Success 200 {object} web.WebResponse{data=[]web.ProductResponse}
// @Failure 500 {object} web.WebResponse
// @Router /products [get]
func (c *ProductController) FindAll(ctx *fiber.Ctx) error {
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
// @Summary Mendapatkan produk berdasarkan ID
// @Description Mencari produk berdasarkan ID produk
// @Security BearerAuth
// @Tags Product
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} web.WebResponse{data=web.ProductResponse}
// @Failure 400,404 {object} web.WebResponse
// @Router /products/{id} [get]
func (c *ProductController) FindById(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid product ID",
		})
	}

	result, err := c.Service.FindById(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Error:  "Product not found",
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   result,
	})
}

// Create godoc
// @Summary Membuat Produk baru
// @Description Membuat produk baru (hanya bisa oleh admin)
// @Tags Product
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body web.ProductCreateOrUpdateRequest true "Product Data"
// @Success 201 {object} web.WebResponse{data=web.ProductResponse}
// @Failure 400 {object} web.WebResponse
// @Router /products [post]
func (c *ProductController) Create(ctx *fiber.Ctx) error {
	var req web.ProductCreateOrUpdateRequest
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
// @Summary Perbarui data produk
// @Description Endpoint ini digunakan untuk memperbarui informasi produk berdasarkan ID.
// @Tags Product
// @Accept json
// @Produce json
// @Param id path int true "ID Produk yang akan diperbarui"
// @Param request body web.ProductCreateOrUpdateRequest true "Data produk terbaru"
// @Success 200 {object} web.WebResponse{data=web.ProductResponse}
// @Failure 400 {object} web.WebResponse
// @Failure 404 {object} web.WebResponse
// @Security BearerAuth
// @Router /products/{id} [put]
func (c *ProductController) Update(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid product ID",
		})
	}

	var req web.ProductCreateOrUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid request body",
		})
	}

	result, err := c.Service.Update(id, req)
	if err != nil {
		if err.Error() == "product not found" {
			return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
				Code:   http.StatusNotFound,
				Status: "NOT FOUND",
				Error:  "Product not found",
			})
		}
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
// @Summary Hapus produk
// @Description Endpoint ini digunakan untuk menghapus produk berdasarkan ID.
// @Tags Product
// @Produce json
// @Param id path int true "ID Produk yang akan dihapus"
// @Success 200 {object} web.WebResponse
// @Failure 400 {object} web.WebResponse
// @Failure 404 {object} web.WebResponse
// @Security BearerAuth
// @Router /products/{id} [delete]
func (c *ProductController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid product ID",
		})
	}

	err = c.Service.Delete(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Error:  "Product not found",
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "Product deleted",
	})
}

// Search godoc
// @Summary Cari, filter, dan paginasi produk
// @Description Endpoint ini digunakan untuk mencari produk berdasarkan kata kunci, menyortir berdasarkan stok, dan menampilkan hasil dengan paginasi.
// @Tags Product
// @Produce json
// @Param q query string false "Kata kunci pencarian (nama produk)"
// @Param sort query string false "Urutkan berdasarkan stok: asc atau desc"
// @Param page query int false "Nomor halaman (default: 1)"
// @Param limit query int false "Jumlah item per halaman (default: 10)"
// @Success 200 {object} web.WebResponse{data=[]web.ProductResponse}
// @Failure 404 {object} web.WebResponse
// @Security BearerAuth
// @Router /products/search [get]
func (c *ProductController) Search(ctx *fiber.Ctx) error {
	keyword := ctx.Query("q", "")
	sort := ctx.Query("sort", "asc")
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	results, err := c.Service.SearchWithFilter(keyword, sort, page, limit)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Error:  err.Error(),
		})
	}

	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   results,
	})
}
