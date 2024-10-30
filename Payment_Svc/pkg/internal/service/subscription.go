package service

import (
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/razorpay/razorpay-go"
	"github.com/spf13/viper"
)

type PaymentService struct {
	Repo        repo.RepoInter
	RazorClient *razorpay.Client
	UserClient  proto.UserServiceClient
	GigClient   proto.GigServiceClient
	proto.UnimplementedPaymentServiceServer
}

func NewPaymentService(repo repo.RepoInter, UserConn proto.UserServiceClient,GigConn proto.GigServiceClient) *PaymentService {
	client := razorpay.NewClient(viper.GetString("ApiKey"), viper.GetString("ApiSecret"))
	return &PaymentService{
		UserClient:  UserConn,
		Repo:        repo,
		GigClient: GigConn,
		RazorClient: client,
	}
}

func (s *PaymentService) CreateSubscription(ctx context.Context, req *proto.CreateSubscriptionRequest) (*proto.CreateSubscriptionResponse, error) {
	userSub, err := s.Repo.GetActiveSubscription(uint(req.UserId))
	if err != nil {
		log.Println("Failed to find sub: ", err.Error())
		return nil, err
	}
	if userSub.ID != 0 && userSub.Active == "Active" {
		return nil, errors.New("user already has an active subscription")
	}

	plan, err := s.Repo.GetPlanByID(uint(req.PlanId))
	if err != nil {
		log.Println("Filed to find the plan: ", err.Error())
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
		log.Println("Failed to create subscription: ", err)
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
		return nil, err
	}

	err = s.Repo.CreatePayment(model.Payment{
		SubscriptionID: subscription["id"].(string),
		UserID:         uint(req.UserId),
		Amount:         plan.Price,
		Status:         "pending",
	})
	if err != nil {
		log.Println("Failed to Save Paymetn Data: ", err.Error())
		return nil, err
	}

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
		log.Printf("Error updating payment for transaction ID %s: %v", transactionID, err)
		return &proto.WebhookResponse{Success: false, Message: "Failed to update payment"}, nil
	}

	subscription, err := s.Repo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		log.Printf("Error fetching subscription for ID %s: %v", subscriptionID, err)
		return &proto.WebhookResponse{Success: false, Message: "Failed to fetch subscription"}, nil
	}

	subscription.Active = "Active"
	err = s.Repo.UpdateSubscription(*subscription)
	if err != nil {
		log.Printf("Error updating subscription for ID %s: %v", subscriptionID, err)
		return &proto.WebhookResponse{Success: false, Message: "Failed to update subscription"}, nil
	}

	log.Printf("Payment captured for transaction ID: %s", transactionID)
	return &proto.WebhookResponse{Success: true, Message: "Payment processed successfully"}, nil
}

func (s *PaymentService) handleSubscriptionCharged(payload map[string]string) (*proto.WebhookResponse, error) {
	subscriptionID := payload["entity.id"]
	amount := payload["entity.amount"]

	subscription, err := s.Repo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		log.Printf("Error fetching subscription for ID %s: %v", subscriptionID, err)
		return &proto.WebhookResponse{Success: false, Message: "Failed to fetch subscription"}, nil
	}

	subscription.Active = "Charged"
	err = s.Repo.UpdateSubscription(*subscription)
	if err != nil {
		log.Printf("Error updating subscription for ID %s: %v", subscriptionID, err)
		return &proto.WebhookResponse{Success: false, Message: "Failed to update subscription"}, nil
	}

	log.Printf("Subscription charged: %s, Amount: %s", subscriptionID, amount)
	return &proto.WebhookResponse{Success: true, Message: "Subscription charged successfully"}, nil
}

func (s *PaymentService) handleSubscriptionCancelled(payload map[string]string) (*proto.WebhookResponse, error) {
	subscriptionID := payload["entity.id"]

	subscription, err := s.Repo.GetSubscriptionByID(subscriptionID)
	if err != nil {
		log.Printf("Error fetching subscription for ID %s: %v", subscriptionID, err)
		return &proto.WebhookResponse{Success: false, Message: "Failed to fetch subscription"}, nil
	}

	subscription.Active = "Cancelled"
	err = s.Repo.UpdateSubscription(*subscription)
	if err != nil {
		log.Printf("Error updating subscription for ID %s: %v", subscriptionID, err)
		return &proto.WebhookResponse{Success: false, Message: "Failed to update subscription"}, nil
	}

	log.Printf("Subscription cancelled: %s", subscriptionID)
	return &proto.WebhookResponse{Success: true, Message: "Subscription cancelled successfully"}, nil
}
