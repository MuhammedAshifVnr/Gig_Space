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

type Order struct {
	gorm.Model
	OrderID      string
	GigID        uint
	ClinetID     uint
	FreelancerID uint
	Status       string
	Amount       int
}

type Quote struct {
	gorm.Model
	GigId        uint64
	Gig          Gig
	ClientId     uint64
	Describe     string
	Price        float64
	DeliveryDays int
}

type CustomGig struct {
	gorm.Model
	GigRequestID uint
	FreelancerID uint
	ClientID     uint
	Title        string
	Description  string
	Price        float64
	DeliveryDays int
}
type CustomOrder struct {
	gorm.Model
	OrderID      string
	CustomGigID  uint
	ClinetID     uint
	FreelancerID uint
	Status       string
	Amount       int
}

type OrderEvent struct{
	OrderID string
	Event string
}