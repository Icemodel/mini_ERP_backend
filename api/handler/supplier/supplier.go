package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/command"
	"mini-erp-backend/api/service/supplier/query"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// CreateSupplier
//
//	@Summary		Create a new supplier
//	@Description	Create a new supplier with the provided information
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			supplier	body		command.CreateSupplierRequest	true	"Supplier information"
//	@Success		201			{object}	model.Supplier
//	@Failure		400			{object}	fiber.Map
//	@Failure		500			{object}	fiber.Map
//	@Router			/api/v1/suppliers [post]
func CreateSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req command.CreateSupplierRequest

		err := c.BodyParser(&req)
		if err != nil {
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

		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

// GetSupplier
//
//	@Summary		Get a supplier by ID
//	@Description	Get supplier details by supplier ID
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Supplier ID (UUID)"
//	@Success		200	{object}	model.Supplier
//	@Failure		400	{object}	fiber.Map
//	@Failure		404	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/suppliers/{id} [get]
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

		result, err := mediatr.Send[*query.GetSupplierRequest, *query.GetSupplierResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to get supplier", "error", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Supplier not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

// UpdateSupplier
//
//	@Summary		Update a supplier
//	@Description	Update supplier information by ID
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			id			path	string							true	"Supplier ID (UUID)"
//	@Param			supplier	body	command.UpdateSupplierRequest	true	"Updated supplier information"
//	@Success		200
//	@Failure		400	{object}	fiber.Map
//	@Failure		404	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/suppliers/{id} [put]
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
		err = c.BodyParser(&req)
		if err != nil {
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

		return c.SendStatus(fiber.StatusOK)
	}
}

// DeleteSupplier
//
//	@Summary		Delete a supplier
//	@Description	Delete a supplier by ID
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Supplier ID (UUID)"
//	@Success		200
//	@Failure		400	{object}	fiber.Map
//	@Failure		404	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/api/v1/suppliers/{id} [delete]
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

		return c.SendStatus(fiber.StatusOK)
	}
}

// GetAllSuppliers
//
//	@Summary		Get all suppliers
//	@Description	Retrieve a list of all suppliers
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			order_by	query	string	false	"Order by field"
//	@Success		200			{array}	model.Supplier
//	@Failure		500			{object}	fiber.Map
//	@Router			/api/v1/suppliers [get]
func GetAllSuppliers(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orderBy := c.Query("order_by", "")

		req := query.GetAllSuppliersRequest{
			OrderBy: orderBy,
		}

		result, err := mediatr.Send[*query.GetAllSuppliersRequest, *query.GetAllSuppliersResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to get all suppliers", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}
