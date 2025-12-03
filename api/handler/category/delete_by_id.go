package category_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/category/command"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// DeleteById is a function to delete category by id
//
//	@Summary		Delete Category by ID
//	@Description	Delete category by ID
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	command.DeleteByIdResult
//	@Failure		404	{object}	map[string]string	"Not Found: Category does not exist"
//	@Failure		500	{object}	map[string]string	"Internal Server Error"
//	@Router			/categories/{id} [delete]
//
//	@param			id	path	string	true	"Category ID"
func DeleteById(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		categoryIdParam := c.Params("id")
		categoryId, err := uuid.Parse(categoryIdParam)
		if err != nil {
			logger.Error("Invalid category ID", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}
		request := command.DeleteByIdRequest{
			CategoryId: categoryId,
		}

		response, err := mediatr.Send[command.DeleteByIdRequest, *command.DeleteByIdResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to delete category", slog.String("error", err.Error()))
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete category",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
