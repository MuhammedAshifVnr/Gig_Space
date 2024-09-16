package db

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(DB string) (*gorm.DB, *redis.Client,*s3.S3) {
	db, err := gorm.Open(postgres.Open(DB), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect DB :%v", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("eu-north-1"),
		Credentials: credentials.NewStaticCredentials(viper.GetString("Accesskey"), viper.GetString("Secretaccesskey"), ""),
	})
	s3Svc := s3.New(sess)
	autoMigrate(db)
	return db, rdb,s3Svc
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.User{},
		&model.Admin{},
		&model.Category{},
		&model.Skills{},
		&model.Freelancer_Skills{},
		&model.UserProfile{},
		&model.ProfilePhoto{},
	)
}
