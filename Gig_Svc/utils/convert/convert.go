package convert

import (
	"bytes"
	"mime/multipart"
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
    fileName := "uploaded_image.jpg"  

    
    file := NewCustomFile(imageBytes, fileName)

    fileHeader := &multipart.FileHeader{
        Filename: fileName,
        Size:     int64(len(imageBytes)),
    }

    return file, fileHeader, nil
}