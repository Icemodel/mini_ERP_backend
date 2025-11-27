package product_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/product/query"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

func Products(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, _ := strconv.Atoi(c.Query("page", "1"))
		pageSize, _ := strconv.Atoi(c.Query("pageSize", "10"))
		search := c.Query("search", "")
		categoryId := c.Query("category_id", "")
		sortBy := c.Query("sortBy", "")
		sortOrder := c.Query("sortOrder", "")

		var parsedCategoryId *uuid.UUID
		if categoryId != "" {
			id, err := uuid.Parse(categoryId)
			if err == nil {
				parsedCategoryId = &id
			}
		}

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

		request := query.ProductsRequest{
			Page:       page,
			PageSize:   pageSize,
			Search:     search,
			CategoryId: parsedCategoryId,
			SortBy:     sortBy,
			SortOrder:  sortOrder,
		}

		response, err := mediatr.Send[query.ProductsRequest, *query.ProductsResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to get products", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get products",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
