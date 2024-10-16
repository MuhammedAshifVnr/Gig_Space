package db

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDB(DB string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(DB), &gorm.Config{})

	if err != nil {
		log.Fatalf("Failed to connect DB :%v", err)
	}
	autoMigrate(db)
	return db
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
