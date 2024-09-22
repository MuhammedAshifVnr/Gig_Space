package helpers

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/model"
)

func GetImageUrls(images []model.Image) string {

	for _, img := range images {
		return img.Url
	}
	return ""
}
