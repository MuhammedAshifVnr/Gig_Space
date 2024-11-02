package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Notificaton_svc/pkg/models"
	"github.com/MuhammedAshifVnr/Gig_Space/Notificaton_svc/utils/helper"
	"github.com/segmentio/kafka-go"
)

func StartRefundConsumer(consumer *kafka.Reader) error {
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		var event models.RefundEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshal payment-done message: %v", err)
			continue
		}
		fmt.Println("Sussess")
		email, err := helper.GetUserEmail(event.UserID)
		if err != nil {
			log.Printf("failed to find email: %v", err)
		}
		sub, Msg := helper.MessageCreater(event.Event, event.OrderID, event.Amoutn)
		if err := helper.SendEmailNotification(email, sub, Msg); err != nil {
			log.Printf("failed to send payment confirmation: %v", err)
		} else {
			log.Printf("payment confirmation sent to user: %v", event.UserID)
		}
	}
}

func StartStatusConsumer(consumer *kafka.Reader) error {
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		var event models.StatusEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshal payment-done message: %v", err)
			continue
		}
		fmt.Println("---")
		userID, err := helper.GetUserID(event.OrderID)
		fmt.Println("user",userID)
		if err != nil {
			log.Printf("failed to find user_id: %v", err)
		}
		email, err := helper.GetUserEmail(uint(userID))
		if err != nil {
			log.Printf("failed to find email: %v", err)
		}
		if err := helper.SendEmailNotification(email, "Order Status Update Notification",
			fmt.Sprintf("Dear Freelancer,\n\nThe status of your order ID %s has been updated to: %s.\n\nBest regards,\nYour Company", event.OrderID, event.Event)); err != nil {
			log.Printf("failed to send payment confirmation: %v", err)
		} else {
			log.Printf("payment confirmation sent to user: %v", userID)
		}
	}
}

func StartPaymentConsumer(consumer *kafka.Reader) error {
	for {
		msg, err := consumer.ReadMessage(context.Background())
		if err != nil {
			return err
		}

		var event models.RefundEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("failed to unmarshal payment-done message: %v", err)
			continue
		}
		fmt.Println("Sussess")
		email, err := helper.GetUserEmail(event.UserID)
		if err != nil {
			log.Printf("failed to find email: %v", err)
		}
		sub, Msg := helper.MessageCreater(event.Event, event.OrderID, event.Amoutn)
		if err := helper.SendEmailNotification(email, sub, Msg); err != nil {
			log.Printf("failed to send payment confirmation: %v", err)
		} else {
			log.Printf("payment confirmation sent to user: %v", event.UserID)
		}
	}
}
