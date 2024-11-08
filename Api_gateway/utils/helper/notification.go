package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)
var redisClient *redis.Client

func SendOflineNotificaton(senderID, recipientID int32) bool {
	ctx := context.Background()
	key := fmt.Sprintf("notify:%d:%d", senderID, recipientID)

	// Try getting the key; if it doesn’t exist, send notification
	exists, err := redisClient.Exists(ctx, key).Result()
	if err != nil {
		fmt.Println("Error accessing Redis:", err)
		return false
	}

	// If key doesn’t exist, notification can be sent
	if exists == 0 {
		// Set key with TTL (e.g., 1 hour)
		redisClient.Set(ctx, key, "1", time.Hour)
		return true
	}

	// Key exists, notification already sent for this sender-recipient pair
	return false
}
