package category_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/category/query"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

type CategoryQuery struct {
	Page      int    `query:"page"`
	PageSize  int    `query:"pageSize"`
	Search    string `query:"search"`
	SortBy    string `query:"sortBy"`
	SortOrder string `query:"sortOrder"`
}

// Categories is a function to get all categories
//
//	@Summary		Get Category list
//	@Description	Get category list
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	query.CategoriesResult
//	@Router			/categories [get]
//
//	@param			page		query	int		false	"Page number"
//	@param			pageSize	query	int		false	"Number of items per page"
//	@param			search		query	string	false	"Search term for name and description"
//	@param			sortBy		query	string	false	"Field to sort by"
//	@param			sortOrder	query	string	false	"Sort order (asc or desc)"
func Categories(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q CategoryQuery

		if err := c.QueryParser(&q); err != nil {
			logger.Error("Failed to parse query parameters", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid query parameters",
			})
		}

		request := query.CategoriesRequest{
			Page:      q.Page,
			PageSize:  q.PageSize,
			Search:    q.Search,
			SortBy:    q.SortBy,
			SortOrder: q.SortOrder,
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