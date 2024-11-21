package upimage

import (
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

func UploadImage(S3 *s3.S3, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	fileName := fileHeader.Filename

	_, err := S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("BUCKET_NAME")),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", viper.GetString("BUCKET_NAME"), fileName)
	return url, nil
}
