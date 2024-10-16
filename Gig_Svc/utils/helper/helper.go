package helper

import (
	"crypto/rand"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/model"
)

func GetImageUrls(images []model.Image) string {

	for _, img := range images {
		return img.Url
	}
	return ""
}

func RandString() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	for i := 0; i < 8; i++ {
		bytes[i] = letters[bytes[i]%byte(len(letters))]
	}
	return string(bytes)
}
