package utils

import "github.com/juanpablocs/video-stream-golang/internal/models"

func ThumbnailMap(thumbnails []models.Thumbnail) map[string]string {
	thumbnailMap := make(map[string]string)
	for _, thumbnail := range thumbnails {
		if thumbnail.Default {
			thumbnailMap[thumbnail.Size] = thumbnail.Path
		}
	}
	return thumbnailMap
}
