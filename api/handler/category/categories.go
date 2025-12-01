package category_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/category/query"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

func Categories(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ตรวจสอบว่ามี query parameters สำหรับ pagination หรือไม่
		hasPageParam := c.Query("page") != ""
		hasPageSizeParam := c.Query("pageSize") != ""
		usePagination := hasPageParam || hasPageSizeParam

		// รับ query parameters
		page, _ := strconv.Atoi(c.Query("page", "1"))
		pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
		search := c.Query("search", "")
		sortBy := c.Query("sortBy", "")
		sortOrder := c.Query("sortOrder", "")

		// จัดการ page ให้ไม่ต่ำกว่า 1
		if page < 1 {
			page = 1
		}

		// จัดการ pageSize
		if pageSize < 1 {
			pageSize = 10 // ค่า default
		}
		if pageSize > 100 {
			pageSize = 100 // จำกัดค่าสูงสุด
		}

		request := query.CategoriesRequest{
			Page:          page,
			PageSize:      pageSize,
			Search:        search,
			SortBy:        sortBy,
			SortOrder:     sortOrder,
			UsePagination: usePagination,
		}

		response, err := mediatr.Send[query.CategoriesRequest, *query.CategoriesResult](c.Context(), request)
		if err != nil {
			logger.Error("Failed to get categories", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get categories",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
