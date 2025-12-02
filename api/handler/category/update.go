package category_handler

import (
	"log/slog"
	"mini-erp-backend/api/service/category/command"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/mehdihadeli/go-mediatr"
)

// Update is a function to update category by id
//
//	@Summary		Update Category by ID
//	@Description	Update category by ID
//	@Tags			Category
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	command.UpdateResult
//	@Router			/categories/{id} [put]
//
//	@param			id			path	string	true	"Category ID"
//	@param			name		body	string	true	"Category Name"
//	@param			description	body	string	false	"Category Description"
func Update(logger *slog.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		categoryIdParam := c.Params("id")
		categoryId, err := uuid.Parse(categoryIdParam)

		if err != nil {
			logger.Error("Invalid category ID", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid category ID",
			})
		}

		request := command.UpdateRequest{
			CategoryId: categoryId,
		}

		if err := c.BodyParser(&request); err != nil {
			logger.Error("Failed to parse update category request", slog.String("error", err.Error()))
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if request.Name == "" {
			logger.Error("Category name is required")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Category name is required",
			})
		}

		response, err := mediatr.Send[command.UpdateRequest, *command.UpdateResult](c.Context(), request)

		if err != nil {
			logger.Error("Failed to update category", slog.String("error", err.Error()))

			// ตรวจสอบว่าเป็น error ชื่อซ้ำหรือไม่
			if strings.Contains(err.Error(), "already exists") {
				return c.Status(fiber.StatusConflict).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			// ตรวจสอบว่าเป็น error ไม่พบข้อมูลหรือไม่
			if strings.Contains(err.Error(), "not found") {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update category",
			})
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
