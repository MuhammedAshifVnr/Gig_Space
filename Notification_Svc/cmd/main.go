package main

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Notificaton_svc/pkg/app"
	"github.com/MuhammedAshifVnr/Gig_Space/Notificaton_svc/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/Notificaton_svc/pkg/di"

	"github.com/spf13/viper"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal("---", err)
	}

	RefundReader := di.NewKafkaConsumer(viper.GetString("Broker"), viper.GetString("RefundTopic"))
	forever := make(chan bool)
	go func() {
		if err := app.StartRefundConsumer(RefundReader); err != nil {
			log.Fatalf("failed to start Refund consumer: %v", err)
		}
	}()

	StatusReader := di.NewKafkaConsumer(viper.GetString("Broker"), viper.GetString("StatusTopic"))
	go func() {
		if err := app.StartStatusConsumer(StatusReader); err != nil {
			log.Fatalf("failed to start Status	 consumer: %v", err)
		}
	}()

	PaymentReader := di.NewKafkaConsumer(viper.GetString("Broker"), viper.GetString("PaymentTopic"))
	go func() {
		if err := app.StartStatusConsumer(PaymentReader); err != nil {
			log.Fatalf("failed to start Payment consumer: %v", err)
		}
	}()

	ForgotReader := di.NewKafkaConsumer(viper.GetString("Broker"), viper.GetString("ForgotTopic"))
	go func() {
		if err := app.StartForgetEmailConsumer(ForgotReader); err != nil {
			log.Fatalf("failed to start Forgot consumer: %v", err)
		}
	}()

	OfflineReader := di.NewKafkaConsumer(viper.GetString("Broker"), viper.GetString("OfflineTopic"))
	go func() {
		if err := app.StartChatNotificationConsumer(OfflineReader); err != nil {
			log.Fatalf("failed to start CharOffline consumer: %v", err)
		}
	}()

	OrderReader := di.NewKafkaConsumer(viper.GetString("Broker"), viper.GetString("OrderTopic"))
	go func() {
		if err := app.StartOrderNotificationConsumer(OrderReader); err != nil {
			log.Fatalf("failed to start Order consumer: %v", err)
		}
	}()

	log.Println("Notification server is running on port ", viper.GetString("Broker"))
	<-forever
}
