package product_handler

import (
	"fmt"
	"log/slog"
	"mini-erp-backend/api/service/product/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

type ProductQuery struct {
	Page       int        `query:"page"`
	PageSize   int        `query:"pageSize"`
	Search     string     `query:"search"`
	CategoryId *uuid.UUID `query:"categoryId"`
	SortBy     string     `query:"sortBy"`
	SortOrder  string     `query:"sortOrder"`
}

func Products(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q ProductQuery

		if err := c.QueryParser(&q); err != nil {
			logger.Error("Failed to parse query parameters", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid query parameters: %v", err),
			})
		}

		request := query.ProductsRequest{
			Page:       q.Page,
			PageSize:   q.PageSize,
			Search:     q.Search,
			CategoryId: q.CategoryId,
			SortBy:     q.SortBy,
			SortOrder:  q.SortOrder,
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
