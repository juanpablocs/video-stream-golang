package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h Handler) VideoByID(c *fiber.Ctx) error {
	videoID := c.Params("id")
	video, err := h.usecase.VideoId(videoID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// TODO: controlled with a flag or parameter
	videoUpload, _ := h.usecase.VideoUploadId(videoID)
	video.VideoUpload = *videoUpload

	return c.JSON(video)
}
