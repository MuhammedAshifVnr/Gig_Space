package config

import (
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

func LoadConfig() error {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	viper.AutomaticEnv()
	return nil
}

func InitKafkaWriters(brokers []string) map[string]*kafka.Writer {
	writers := make(map[string]*kafka.Writer)
	topics := []string{viper.GetString("FORGOT_TOPIC"), viper.GetString("STATUS_TOPIC"),viper.GetString("ORDER_TOPIC")}
	for _, topic := range topics {
		writers[topic] = NewKafkaWriter(brokers, topic)
	}
	return writers
}

func NewKafkaWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		BatchTimeout: 10 * time.Millisecond,
	}
}
