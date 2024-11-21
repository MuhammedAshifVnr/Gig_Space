package upload

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

func UploadPhoto(S3 *s3.S3, Pic []byte, UserId uint) (string, error) {
	PhotoBytes := Pic
	fileName := "profile_" + string(rune(UserId)) + ".jpg"
	_, err := S3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(viper.GetString("BUCKET_NAME")),
		Key:    aws.String(fileName),
		Body:   bytes.NewReader(PhotoBytes),
	})
	if err != nil {
		return "", err
	}
	photoURL := "https://" + viper.GetString("BUCKET_NAME") + ".s3.amazonaws.com/" + fileName
	return photoURL, nil
}
