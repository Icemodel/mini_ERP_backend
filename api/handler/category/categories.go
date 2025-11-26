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
		// รับ query parameters
		page, _ := strconv.Atoi(c.Query("page", "0"))
		pageSize, _ := strconv.Atoi(c.Query("pageSize", "0"))
		search := c.Query("search", "")

		request := query.CategoriesRequest{
			Page:     page,
			PageSize: pageSize,
			Search:   search,
		}

		result, err := mediatr.Send[query.CategoriesRequest, *query.CategoriesResult](c.Context(), request)
		if err != nil {
			logger.Error("Failed to get categories", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get categories",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
