package controller

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"inventory-management-api/model/web"
	"inventory-management-api/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type StockMovementController struct {
	Service service.StockMovementService
}

func NewStockMovementController(s service.StockMovementService) *StockMovementController {
	return &StockMovementController{Service: s}
}

// FindAll godoc
// @Summary Ambil semua data pergerakan stok
// @Description Endpoint ini digunakan untuk mengambil seluruh data pergerakan stok (masuk & keluar).
// @Tags StockMovement
// @Produce json
// @Success 200 {object} web.WebResponse{data=[]web.StockMovementResponse}
// @Failure 500 {object} web.WebResponse
// @Security BearerAuth
// @Router /stock-movements [get]
func (c *StockMovementController) FindAll(ctx *fiber.Ctx) error {
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
// @Summary Ambil data pergerakan stok berdasarkan ID
// @Description Endpoint ini mengambil satu data pergerakan stok berdasarkan ID-nya.
// @Tags StockMovement
// @Produce json
// @Param id path int true "ID pergerakan stok"
// @Success 200 {object} web.WebResponse{data=web.StockMovementResponse}
// @Failure 400,404,500 {object} web.WebResponse
// @Security BearerAuth
// @Router /stock-movements/{id} [get]
func (c *StockMovementController) FindById(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid ID",
		})
	}

	result, err := c.Service.FindById(id)
	if err != nil {
		if err.Error() == "stock movement not found" {
			return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
				Code:   http.StatusNotFound,
				Status: "NOT FOUND",
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
		Data:   result,
	})
}

// Create godoc
// @Summary Tambah data pergerakan stok baru
// @Description Endpoint ini digunakan untuk menambah pergerakan stok masuk atau keluar.
// @Tags StockMovement
// @Accept json
// @Produce json
// @Param request body web.StockMovementCreateRequest true "Data pergerakan stok"
// @Success 201 {object} web.WebResponse{data=web.StockMovementResponse}
// @Failure 400,401,500 {object} web.WebResponse
// @Security BearerAuth
// @Router /stock-movements [post]
func (c *StockMovementController) Create(ctx *fiber.Ctx) error {
	var req web.StockMovementCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid request body",
		})
	}

	userIDVal := ctx.Locals("user_id")
	userID, ok := userIDVal.(int)
	if !ok || userID == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(web.WebResponse{
			Code:   http.StatusUnauthorized,
			Status: "UNAUTHORIZED",
			Error:  "User ID not found in token",
		})
	}

	result, err := c.Service.Create(userID, req)
	if err != nil {
		msg := err.Error()
		switch msg {
		case "validation failed", "product not found", "stock not enough":
			return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "BAD REQUEST",
				Error:  msg,
			})
		default:
			return ctx.Status(http.StatusInternalServerError).JSON(web.WebResponse{
				Code:   http.StatusInternalServerError,
				Status: "INTERNAL SERVER ERROR",
				Error:  msg,
			})
		}
	}

	return ctx.Status(http.StatusCreated).JSON(web.WebResponse{
		Code:   http.StatusCreated,
		Status: "CREATED",
		Data:   result,
	})
}

// Delete godoc
// @Summary Hapus data pergerakan stok
// @Description Endpoint ini digunakan untuk menghapus pergerakan stok berdasarkan ID.
// @Tags StockMovement
// @Produce json
// @Param id path int true "ID pergerakan stok"
// @Success 200 {object} web.WebResponse{data=string}
// @Failure 400,404,500 {object} web.WebResponse
// @Security BearerAuth
// @Router /stock-movements/{id} [delete]
func (c *StockMovementController) Delete(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
			Code:   http.StatusBadRequest,
			Status: "BAD REQUEST",
			Error:  "Invalid ID",
		})
	}

	err = c.Service.Delete(id)
	if err != nil {
		if err.Error() == "stock movement not found" {
			return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
				Code:   http.StatusNotFound,
				Status: "NOT FOUND",
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
		Data:   "Stock movement deleted",
	})
}

// GetMonthlyReport godoc
// @Summary Ambil laporan stok bulanan
// @Description Mengambil laporan pergerakan stok berdasarkan bulan, bisa difilter berdasarkan user, produk, atau tipe, dan bisa diekspor ke CSV.
// @Tags StockMovement
// @Produce json
// @Param month query string false "Format bulan: YYYY-MM (contoh: 2024-06)"
// @Param user_id query int false "Filter berdasarkan ID user"
// @Param product_id query int false "Filter berdasarkan ID produk"
// @Param type query string false "Jenis pergerakan (in atau out)"
// @Param export query string false "Jika bernilai 'csv', maka file akan didownload dalam format CSV"
// @Success 200 {object} web.WebResponse{data=[]web.StockMovementResponse}
// @Failure 400,404,500 {object} web.WebResponse
// @Security BearerAuth
// @Router /reports/stock-movements [get]
func (c *StockMovementController) GetMonthlyReport(ctx *fiber.Ctx) error {
	month := ctx.Query("month")
	export := ctx.Query("export")
	userID := ctx.Query("user_id")
	productID := ctx.Query("product_id")
	movementType := ctx.Query("type")

	filters := map[string]interface{}{}

	// Validasi dan parsing user_id
	if userID != "" {
		if id, err := strconv.Atoi(userID); err == nil {
			filters["user_id"] = id
		} else {
			return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "BAD REQUEST",
				Error:  "invalid user_id",
			})
		}
	}

	// Validasi dan parsing product_id
	if productID != "" {
		if id, err := strconv.Atoi(productID); err == nil {
			filters["product_id"] = id
		} else {
			return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "BAD REQUEST",
				Error:  "invalid product_id",
			})
		}
	}

	// Validasi type hanya boleh "in" atau "out"
	if movementType != "" {
		if movementType != "in" && movementType != "out" {
			return ctx.Status(http.StatusBadRequest).JSON(web.WebResponse{
				Code:   http.StatusBadRequest,
				Status: "BAD REQUEST",
				Error:  "invalid movement type, must be 'in' or 'out'",
			})
		}
		filters["type"] = movementType
	}

	// Ambil data dari service
	data, err := c.Service.GetMonthlyReport(month, filters)
	if err != nil {
		if err.Error() == "report not found" {
			return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
				Code:   http.StatusNotFound,
				Status: "NOT FOUND",
				Error:  "Report not found or filter returned no data",
			})
		}
		return ctx.Status(http.StatusInternalServerError).JSON(web.WebResponse{
			Code:   http.StatusInternalServerError,
			Status: "INTERNAL SERVER ERROR",
			Error:  err.Error(),
		})
	}

	// Jika hasil kosong, balas 404
	if len(data) == 0 {
		return ctx.Status(http.StatusNotFound).JSON(web.WebResponse{
			Code:   http.StatusNotFound,
			Status: "NOT FOUND",
			Error:  "No stock movement data found for the given filter",
		})
	}

	// Export CSV jika diminta
	if export == "csv" {
		var b bytes.Buffer
		writer := csv.NewWriter(&b)

		writer.Write([]string{"ID", "Product ID", "User ID", "Type", "Quantity", "Note", "Created At"})
		for _, m := range data {
			writer.Write([]string{
				strconv.Itoa(m.ID),
				strconv.Itoa(m.ProductID),
				strconv.Itoa(m.UserID),
				m.Type,
				strconv.Itoa(m.Quantity),
				m.Note,
				m.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}
		writer.Flush()

		timestamp := time.Now().Format("20060102_150405")
		filename := fmt.Sprintf("stock_report_%s.csv", timestamp)

		ctx.Set("Content-Type", "text/csv")
		ctx.Set("Content-Disposition", "attachment; filename="+filename)
		return ctx.Send(b.Bytes())
	}

	// Response sukses
	return ctx.JSON(web.WebResponse{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   data,
	})
}
