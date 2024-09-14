package model

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Gig struct {
	gorm.Model
	Title        string
	Description  string
	Category     uint
	FrelancerID  uint
	Price        float64
	DeliveryDays int
	Revisions    int
	Image_Urls   pq.StringArray `gorm:"type:text[]"`
}
