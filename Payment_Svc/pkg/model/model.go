package model

import (
	"time"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	Name     string
	Price    int
	RazorpayPlanID string 
	Period         string 
	Interval       int
	Description    string  
}

type Subscription struct {
	gorm.Model
	UserID         uint
	SubscriptionID string
	OrderID        string
	Active         string
	TrialUsed      bool
	StartDate      time.Time
	EndDate        time.Time
}

type Payment struct {
	gorm.Model
	SubscriptionID string
	UserID         uint
	Amount         int
	Status         string
	TransactionID  string
}
