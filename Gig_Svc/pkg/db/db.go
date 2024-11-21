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
	DB+="gig_svc"
	db, err := gorm.Open(postgres.Open(DB), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB :%v", err)
	}
	autoMigrate(db)
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(viper.GetString("ACCESS_KEY"), viper.GetString("SECRET_ACCESS_KEY"), ""),
	})
	s3Svc := s3.New(sess)

	return db, s3Svc
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Gig{},
		&model.Image{},
		&model.Order{},
		&model.Quote{},
		&model.CustomGig{},
		&model.CustomOrder{},
	)
}
