package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/command"
	"mini-erp-backend/api/service/supplier/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// CreateSupplier godoc
// @Summary Create a new supplier
// @Description Create a new supplier with the provided information
// @Tags suppliers
// @Accept json
// @Produce json
// @Param supplier body command.CreateSupplierRequest true "Supplier information"
// @Success 201 {object} model.Supplier
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /suppliers [post]
func CreateSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req command.CreateSupplierRequest
		
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		result, err := mediatr.Send[*command.CreateSupplierRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to create supplier", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Supplier created successfully",
			"data":    result,
		})
	}
}

// GetSupplier godoc
// @Summary Get a supplier by ID
// @Description Get supplier details by supplier ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID (UUID)"
// @Success 200 {object} model.Supplier
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /suppliers/{id} [get]
func GetSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		supplierId, err := uuid.Parse(idParam)
		if err != nil {
			logger.Error("Invalid supplier ID", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid supplier ID",
			})
		}

		req := query.GetSupplierRequest{
			SupplierId: supplierId,
		}

		result, err := mediatr.Send[*query.GetSupplierRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to get supplier", "error", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Supplier not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Supplier retrieved successfully",
			"data":    result,
		})
	}
}

// UpdateSupplier godoc
// @Summary Update a supplier
// @Description Update supplier information by ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID (UUID)"
// @Param supplier body command.UpdateSupplierRequest true "Updated supplier information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /suppliers/{id} [put]
func UpdateSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		supplierId, err := uuid.Parse(idParam)
		if err != nil {
			logger.Error("Invalid supplier ID", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid supplier ID",
			})
		}

		var req command.UpdateSupplierRequest
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		req.SupplierId = supplierId

		_, err = mediatr.Send[*command.UpdateSupplierRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to update supplier", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Supplier updated successfully",
		})
	}
}

// DeleteSupplier godoc
// @Summary Delete a supplier
// @Description Delete a supplier by ID
// @Tags suppliers
// @Accept json
// @Produce json
// @Param id path string true "Supplier ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /suppliers/{id} [delete]
func DeleteSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		supplierId, err := uuid.Parse(idParam)
		if err != nil {
			logger.Error("Invalid supplier ID", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid supplier ID",
			})
		}

		req := command.DeleteSupplierRequest{
			SupplierId: supplierId,
		}

		_, err = mediatr.Send[*command.DeleteSupplierRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to delete supplier", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Supplier deleted successfully",
		})
	}
}

// GetAllSuppliers godoc
// @Summary Get all suppliers
// @Description Retrieve a list of all suppliers
// @Tags suppliers
// @Accept json
// @Produce json
// @Param order_by query string false "Order by field (default: created_at DESC)"
// @Success 200 {array} model.Supplier
// @Failure 500 {object} map[string]interface{}
// @Router /suppliers [get]
func GetAllSuppliers(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orderBy := c.Query("order_by", "")

		req := query.GetAllSuppliersRequest{
			OrderBy: orderBy,
		}

		result, err := mediatr.Send[*query.GetAllSuppliersRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to get all suppliers", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Suppliers retrieved successfully",
			"data":    result,
		})
	}
}

// SearchSuppliers godoc
// @Summary Search suppliers
// @Description Search suppliers by email and/or name
// @Tags suppliers
// @Accept json
// @Produce json
// @Param email query string false "Filter by email"
// @Param name query string false "Filter by name"
// @Param order_by query string false "Order by field (default: name ASC)"
// @Success 200 {array} model.Supplier
// @Failure 500 {object} map[string]interface{}
// @Router /suppliers/search [get]
func SearchSuppliers(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		email := c.Query("email", "")
		name := c.Query("name", "")
		orderBy := c.Query("order_by", "")

		req := query.SearchSuppliersRequest{
			Email:   email,
			Name:    name,
			OrderBy: orderBy,
		}

		result, err := mediatr.Send[*query.SearchSuppliersRequest, interface{}](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to search suppliers", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Suppliers search completed successfully",
			"data":    result,
		})
	}
}
