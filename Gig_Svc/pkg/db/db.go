package db

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(DB string) (*gorm.DB, *s3.S3) {
	db, err := gorm.Open(postgres.Open(DB), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB :%v", err)
	}
	autoMigrate(db)
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(viper.GetString("Accesskey"), viper.GetString("Secretaccesskey"), ""),
	})
	s3Svc := s3.New(sess)

	return db, s3Svc
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Gig{},
		&model.Image{},
	)
}
