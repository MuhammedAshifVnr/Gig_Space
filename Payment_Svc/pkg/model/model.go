package model

import (
	"time"

	"gorm.io/gorm"
)

type Plan struct {
	gorm.Model
	Name           string
	Price          int
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

type OrderPayment struct {
	gorm.Model
	ReceiptID     string
	OrderID       string
	TransactionID string
	Status        string
	UserID        uint
	Amount        int
}

type Wallet struct {
	gorm.Model
	UserID              uint `gorm:"unique"`
	Balance             int64
	Pin_hash            string
	Fund_account_id     string
	Last_transaction_at time.Time
}

type ForgotEvent struct{
	Otp string
	Email string
	Event string
}

type StatusEvent struct{
	OrderID string
	User_id uint
	Event   string
}