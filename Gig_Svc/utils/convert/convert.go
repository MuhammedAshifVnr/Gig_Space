package convert

import (
	"bytes"
	"math/rand"
	"mime/multipart"
	"time"
)

type CustomFile struct {
	*bytes.Reader
	Name string
}

func NewCustomFile(data []byte, name string) *CustomFile {
	return &CustomFile{
		Reader: bytes.NewReader(data),
		Name:   name,
	}
}

func (cf *CustomFile) Close() error {
	return nil
}

func ConvertToMultipartFile(imageBytes []byte) (multipart.File, *multipart.FileHeader, error) {
	fileName := generateRandomString(6)

	file := NewCustomFile(imageBytes, fileName)

	fileHeader := &multipart.FileHeader{
		Filename: fileName,
		Size:     int64(len(imageBytes)),
	}

	return file, fileHeader, nil
}

func generateRandomString(length int) string {

	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	result := make([]byte, length)

	for i := range result {
		result[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(result)
}
