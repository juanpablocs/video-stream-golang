package utils

import "github.com/juanpablocs/video-stream-golang/internal/models"

func FilterResolutions(videoWidth, videoHeight int) []models.Resolution {
	resolutions := []models.Resolution{
		{Width: 1920, Height: 1080, Label: "1080p"},
		{Width: 1280, Height: 720, Label: "720p"},
		{Width: 854, Height: 480, Label: "480p"},
		{Width: 640, Height: 360, Label: "360p"},
		{Width: 426, Height: 240, Label: "240p"},
	}

	var applicableResolutions []models.Resolution
	for _, res := range resolutions {
		if videoWidth >= res.Width && videoHeight >= res.Height {
			applicableResolutions = append(applicableResolutions, res)
		}
	}
	// Si el video es de menor resolución que todas las disponibles,
	// se agrega solo la resolución más baja a aplicables.
	if len(applicableResolutions) == 0 && len(resolutions) > 0 {
		lowestResolution := resolutions[len(resolutions)-1]
		// Verifica si la resolución del video es igual o menor a la más baja
		if videoWidth <= lowestResolution.Width && videoHeight <= lowestResolution.Height {
			applicableResolutions = append(applicableResolutions, lowestResolution)
		}
	}
	return applicableResolutions
}
