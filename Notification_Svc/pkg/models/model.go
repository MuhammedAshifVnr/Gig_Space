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