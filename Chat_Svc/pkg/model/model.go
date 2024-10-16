package model

import "time"

type Message struct {
	SenderID    int32
	RecipientID int32
	MessageText string
	CreatedAt   time.Time
}
