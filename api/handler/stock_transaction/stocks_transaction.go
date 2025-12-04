package stocktransaction_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/stock_transaction/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

type StockTransactionQuery struct {
	Page      int        `query:"page"`
	PageSize  int        `query:"pageSize"`
	Search    string     `query:"search"`
	ProductId *uuid.UUID `query:"productId"`
	SortBy    string     `query:"sortBy"`
	SortOrder string     `query:"sortOrder"`
}

// StockTransactions is a function to get all stock transactions
//
//	@Summary		Get Stock Transactions list
//	@Description	Get stock transactions list
//	@Tags				StockTransaction
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	query.StocksResult
//	@Router			/stocks [get]
//
//	@param			page		query	int		false	"Page number"
//	@param			pageSize	query	int		false	"Number of items per page"
//	@param			search		query	string	false	"Search term for quantity, type, and reason"
//	@param			productId	query	string	false	"Filter by Product ID"
//	@param			sortBy		query	string	false	"Field to sort by"
//	@param			sortOrder	query	string	false	"Sort order (asc or desc)"
func StockTransactions(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var q StockTransactionQuery

		if err := c.QueryParser(&q); err != nil {
			logger.Error("Failed to parse query parameters", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid query parameters",
			})
		}

		request := query.StocksRequest{
			Page:      q.Page,
			PageSize:  q.PageSize,
			Search:    q.Search,
			ProductId: q.ProductId,
			SortBy:    q.SortBy,
			SortOrder: q.SortOrder,
		}

		response, err := mediatr.Send[query.StocksRequest, *query.StocksResult](c.Context(), request)
		if err != nil {
			logger.Error("Failed to get stock transactions", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to get stock transactions",
			})
		}
		return c.Status(fiber.StatusOK).JSON(response)
	}
}
