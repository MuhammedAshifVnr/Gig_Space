package helper

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func SendNotification(ctx context.Context, event interface{}, key string, writer *kafka.Writer) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(key),
			Value: eventData,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}
	return nil
}
