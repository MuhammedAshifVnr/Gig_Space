package db

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(DB string) (*gorm.DB, *redis.Client) {
	DB+="payment_svc"
	db, err := gorm.Open(postgres.Open(DB), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect DB :%v", err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	autoMigrate(db)
	return db, rdb
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&model.Subscription{},
		&model.Payment{},
		&model.Plan{},
		&model.OrderPayment{},
		&model.Wallet{},
	)
}
