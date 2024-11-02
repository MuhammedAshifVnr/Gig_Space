package config

import (
	"time"

	"github.com/segmentio/kafka-go"
)

func NewKafkaWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...), // e.g., "localhost:9092"
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond, 
	}
}

func InitKafkaWriters(brokers []string, topics []string) map[string]*kafka.Writer {
	writers := make(map[string]*kafka.Writer)
	for _, topic := range topics {
		writers[topic] = NewKafkaWriter(brokers, topic)
	}
	return writers
}