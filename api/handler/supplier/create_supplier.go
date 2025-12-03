package supplier

import (
	"log/slog"
	"mini-erp-backend/api/service/supplier/command"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/mehdihadeli/go-mediatr"
)

// CreateSupplier
//
//	@Summary		Create a new supplier
//	@Description	Create a new supplier with the provided information
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Param			supplier	body	command.CreateSupplierRequest	true	"Supplier information"
//	@Success		201	{object}	model.Supplier
//	@Failure		400	{object}	api.ErrorResponse
//	@Failure		500	{object}	api.ErrorResponse
//	@Router			/suppliers [post]
var phoneRegex = regexp.MustCompile(`^[\d\s\-\+\(\)]+$`)

func CreateSupplier(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req command.CreateSupplierRequest
		
		err := c.BodyParser(&req)
		if err != nil {
			logger.Error("Failed to parse request body", "error", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Validate phone number format
		phoneRegex := regexp.MustCompile(`^[\d\s\-\+\(\)]+$`)
		if !phoneRegex.MatchString(req.Phone) {
			logger.Error("Invalid phone number format", "phone", req.Phone)
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid phone number format"})
        }

		// Validate email format
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(req.Email) {
			logger.Error("Invalid email format", "email", req.Email)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
		}
		
		result, err := mediatr.Send[*command.CreateSupplierRequest, *command.CreateSupplierResult](c.Context(), &req)
		if err != nil {
			logger.Error("Failed to create supplier", "error", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		
		return c.Status(fiber.StatusCreated).JSON(result)
	}
}

