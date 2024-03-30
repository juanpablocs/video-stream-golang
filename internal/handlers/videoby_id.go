package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h Handler) VideoByID(c *fiber.Ctx) error {
	video, err := h.usecase.VideoId(c.Params("id"))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(video)
}
