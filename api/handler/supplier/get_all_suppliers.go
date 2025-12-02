package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/query"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// AllSuppliers
//
// 	@Summary		Get all suppliers
// 	@Description	Retrieve a list of all suppliers
// 	@Tags			Supplier
// 	@Accept			json
// 	@Produce		json
// 	@Param			order_by	query	string	false	"Order by field"
// 	@Success		200	{array}	model.Supplier
// 	@Failure		500	{object}	api.ErrorResponse
// 	@Router			/suppliers [get]
func AllSuppliers(logger *slog.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        orderBy := c.Query("order_by", "")
        req := query.AllSuppliersRequest{OrderBy: orderBy}

        result, err := mediatr.Send[*query.AllSuppliersRequest, *query.AllSuppliersResult](c.Context(), &req)
        if err != nil {
            logger.Error("Failed to get all suppliers", "error", err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
        }

        return c.Status(fiber.StatusOK).JSON(result)
    }
}
