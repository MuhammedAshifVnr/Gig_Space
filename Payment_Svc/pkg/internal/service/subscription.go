package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/logger"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/razorpay/razorpay-go"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type PaymentService struct {
	Repo        repo.RepoInter
	RazorClient *razorpay.Client
	UserClient  proto.UserServiceClient
	GigClient   proto.GigServiceClient
	kafkaWriter map[string]*kafka.Writer
	Log         *logrus.Logger
	proto.UnimplementedPaymentServiceServer
}

func NewPaymentService(repo repo.RepoInter, UserConn proto.UserServiceClient, GigConn proto.GigServiceClient, kafkaWriter map[string]*kafka.Writer) *PaymentService {
	client := razorpay.NewClient(viper.GetString("API_KEY"), viper.GetString("API_SECRET"))
	return &PaymentService{
		UserClient:  UserConn,
		Repo:        repo,
		GigClient:   GigConn,
		kafkaWriter: kafkaWriter,
		RazorClient: client,
		Log:         logger.Log,
	}
}

func (s *PaymentService) CreateSubscription(ctx context.Context, req *proto.CreateSubscriptionRequest) (*proto.CreateSubscriptionResponse, error) {
	userSub, err := s.Repo.GetActiveSubscription(uint(req.UserId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
		}).Error("Failed to find subscription: ", err)
		return nil, err
	}
	if userSub.ID != 0 && userSub.Active == "Active" {
		return nil, errors.New("user already has an active subscription")
	}

	plan, err := s.Repo.GetPlanByID(uint(req.PlanId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"plan_id": req.PlanId,
		}).Error("Failed to find the plan: ", err)
		return nil, err
	}

	subscriptionData := map[string]interface{}{
		"plan_id":         plan.RazorpayPlanID,
		"total_count":     12,
		"quantity":        1,
		"customer_notify": 1,
		"expire_by":       time.Now().AddDate(1, 0, 0).Unix(),
		"notes": map[string]interface{}{
			"user_id": req.UserId,
		},
	}

	subscription, err := s.RazorClient.Subscription.Create(subscriptionData, nil)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"plan_id": plan.RazorpayPlanID,
			"user_id": req.UserId,
		}).Error("Failed to create subscription: ", err)
		return nil, err
	}

	err = s.Repo.CreateSubscription(model.Subscription{
		UserID:         uint(req.UserId),
		SubscriptionID: subscription["id"].(string),
		Active:         "InActive",
		StartDate:      time.Now(),
		EndDate:        time.Now().AddDate(0, 1, 0),
	})
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id":         req.UserId,
			"subscription_id": subscription["id"].(string),
		}).Error("Failed to save subscription in the database: ", err)
		return nil, err
	}

	err = s.Repo.CreatePayment(model.Payment{
		SubscriptionID: subscription["id"].(string),
		UserID:         uint(req.UserId),
		Amount:         plan.Price,
		Status:         "pending",
	})
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id":         req.UserId,
			"subscription_id": subscription["id"].(string),
		}).Error("Failed to save payment data: ", err)
		return nil, err
	}

	s.Log.WithFields(logrus.Fields{
		"user_id":         req.UserId,
		"subscription_id": subscription["id"].(string),
	}).Info("Subscription created and payment data saved successfully")

	return &proto.CreateSubscriptionResponse{
		Message: "Payment link : " + subscription["short_url"].(string),
		Success: true,
	}, nil
}

func (s *PaymentService) HandleWebhook(ctx context.Context, req *proto.WebhookRequest) (*proto.WebhookResponse, error) {
	eventType := req.Payload["event"]

	switch eventType {
	case "payment.captured":
		return s.handlePaymentCaptured(req.Payload)
	case "subscription.charged":
		return s.handleSubscriptionCharged(req.Payload)
	case "subscription.cancelled":
		return s.handleSubscriptionCancelled(req.Payload)
	default:
		log.Println("Unhandled event type: ", eventType)
		return &proto.WebhookResponse{Success: false, Message: "Unhandled event"}, nil
	}
}

func (s *PaymentService) handlePaymentCaptured(payload map[string]string) (*proto.WebhookResponse, error) {
	transactionID := payload["entity.id"]
	subscriptionID := payload["entity.subscription_id"]
	amount, _ := strconv.Atoi(payload["entity.amount"])

	payment := model.Payment{
		TransactionID:  transactionID,
		SubscriptionID: subscriptionID,
		Status:         "completed",
		Amount:         amount,
	}

	err := s.Repo.UpdatePayment(payment)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"transaction_id": transactionID,
			"error":          err,
		}).Error("Failed to update payment")
		return &proto.WebhookResponse{Success: false, Message: "Failed to update payment"}, nil
	}

	subscription, err := s.Repo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"subscription_id": subscriptionID,
			"error":           err,
		}).Error("Failed to fetch subscription")
		return &proto.WebhookResponse{Success: false, Message: "Failed to fetch subscription"}, nil
	}

	subscription.Active = "Active"
	err = s.Repo.UpdateSubscription(*subscription)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"subscription_id": subscriptionID,
			"error":           err,
		}).Error("Failed to update subscription")
		return &proto.WebhookResponse{Success: false, Message: "Failed to update subscription"}, nil
	}

	s.Log.WithFields(logrus.Fields{
		"transaction_id": transactionID,
	}).Info("Payment captured and subscription updated successfully")

	return &proto.WebhookResponse{Success: true, Message: "Payment processed successfully"}, nil
}

func (s *PaymentService) handleSubscriptionCharged(payload map[string]string) (*proto.WebhookResponse, error) {
	subscriptionID := payload["entity.id"]
	amount := payload["entity.amount"]

	subscription, err := s.Repo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"subscription_id": subscriptionID,
			"error":           err,
		}).Error("Failed to fetch subscription")
		return &proto.WebhookResponse{Success: false, Message: "Failed to fetch subscription"}, nil
	}

	subscription.Active = "Charged"
	err = s.Repo.UpdateSubscription(*subscription)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"subscription_id": subscriptionID,
			"error":           err,
		}).Error("Failed to update subscription")
		return &proto.WebhookResponse{Success: false, Message: "Failed to update subscription"}, nil
	}

	s.Log.WithFields(logrus.Fields{
		"subscription_id": subscriptionID,
		"amount":          amount,
	}).Info("Subscription charged successfully")
	return &proto.WebhookResponse{Success: true, Message: "Subscription charged successfully"}, nil
}

func (s *PaymentService) handleSubscriptionCancelled(payload map[string]string) (*proto.WebhookResponse, error) {
	subscriptionID := payload["entity.id"]

	subscription, err := s.Repo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"subscription_id": subscriptionID,
			"error":           err,
		}).Error("Failed to fetch subscription")
		return &proto.WebhookResponse{Success: false, Message: "Failed to fetch subscription"}, nil
	}

	subscription.Active = "Cancelled"
	err = s.Repo.UpdateSubscription(*subscription)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"subscription_id": subscriptionID,
			"error":           err,
		}).Error("Failed to update subscription")
		return &proto.WebhookResponse{Success: false, Message: "Failed to update subscription"}, nil
	}

	s.Log.WithFields(logrus.Fields{
		"subscription_id": subscriptionID,
	}).Info("Subscription cancelled successfully")
	return &proto.WebhookResponse{Success: true, Message: "Subscription cancelled successfully"}, nil
}
