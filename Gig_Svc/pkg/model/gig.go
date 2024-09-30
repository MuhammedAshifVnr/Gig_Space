package model

import (
	"gorm.io/gorm"
)

type Gig struct {
	gorm.Model
	Title        string
	Description  string
	Category     string
	FreelancerID uint
	Price        float64
	DeliveryDays int
	Revisions    int
	Images       []Image `gorm:"foreignKey:GigID"` // Relationship to Images
}

type Image struct {
	gorm.Model
	GigID uint
	Url   string
}
