package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h Handler) PublicListVideos(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	pageSize := c.Query("pageSize", "10")
	pageInt, _ := strconv.ParseInt(page, 10, 64)
	pageSizeInt, _ := strconv.ParseInt(pageSize, 10, 64)
	data, err := h.usecase.VideoList(pageInt, pageSizeInt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(data)
}
