package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/query"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// GetAllSuppliers
//
// 	@Summary		Get all suppliers
// 	@Description	Retrieve a list of all suppliers
// 	@Tags			Supplier
// 	@Accept			json
// 	@Produce		json
// 	@Param			order_by	query	string	false	"Order by field"
// 	@Success		200	{array}	model.Supplier
// 	@Failure		500	{object}	fiber.Map
// 	@Router			/api/v1/suppliers [get]
func GetAllSuppliers(logger *slog.Logger) fiber.Handler {
    return func(c *fiber.Ctx) error {
        orderBy := c.Query("order_by", "")
        req := query.GetAllSuppliersRequest{OrderBy: orderBy}

        result, err := mediatr.Send[*query.GetAllSuppliersRequest, *query.GetAllSuppliersResult](c.Context(), &req)
        if err != nil {
            logger.Error("Failed to get all suppliers", "error", err)
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
        }

        return c.Status(fiber.StatusOK).JSON(result)
    }
}
