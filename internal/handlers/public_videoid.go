package handlers

import "github.com/gofiber/fiber/v2"

func (h Handler) PublicVideoID(c *fiber.Ctx) error {
	ID := c.Params("id")
	data, err := h.usecase.VideoIdPlayer(ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(data)
}
