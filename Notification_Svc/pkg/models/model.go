package models

type RefundEvent struct {
	UserID  uint `json:"user_id"`
	OrderID string
	Event   string
	Amoutn  int
}

type StatusEvent struct{
	OrderID string
	Event   string
}

type ForgotEvent struct{
	Otp string
	Email string
	Event string
}

type ChatEvent struct{
	SenderID    int32
	RecipientID int32
	Event       string
}