package di

import (
	"github.com/segmentio/kafka-go"
)

func NewKafkaConsumer(brokers string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokers},
		GroupID: "notification-service",
		Topic:   topic,
	})
}
