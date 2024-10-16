package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"sync"
// 	"time"

// 	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/internal/repo"
// 	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
// 	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/utils/helper"
// 	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/utils/payment"
// 	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
// 	"github.com/razorpay/razorpay-go"
// 	"github.com/spf13/viper"
// )

// type PaymentService struct {
// 	Repo        repo.RepoInter
// 	RazorClient *razorpay.Client
// 	proto.UnimplementedPaymentServiceServer
// }

// func NewPaymentService(repo repo.RepoInter) *PaymentService {
// 	client := razorpay.NewClient(viper.GetString("ApiKey"), viper.GetString("ApiSecret"))
// 	return &PaymentService{
// 		Repo:        repo,
// 		RazorClient: client,
// 	}
// }

// func (s *PaymentService) CreateSubscription(ctx context.Context, req *proto.CreateSubscriptionRequest) (*proto.CreateSubscriptionResponse, error) {
// 	var wg sync.WaitGroup
// 	var receipt string

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		receipt = helper.RandString()
// 	}()

// 	userSub, err := s.Repo.GetActiveSubscription(uint(req.UserId))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if userSub.ID != 0 && userSub.Active == "Active" {
// 		return nil, errors.New("user already has an active subscription")
// 	}

// 	plan, err := s.Repo.GetPlanByID(uint(req.PlanId))
// 	if err != nil {
// 		return nil, err
// 	}
// 	wg.Wait()

// 	orderData := map[string]interface{}{
// 		"amount":   plan.Price * 100,
// 		"currency": "INR",
// 		"receipt":  receipt,
// 	}
// 	order, err := s.RazorClient.Order.Create(orderData, nil)
// 	if err != nil {
// 		log.Println("Faild to create order : ", err)
// 		return nil, err
// 	}
// 	startDate := time.Now()
// 	err = s.Repo.CreateSubscription(model.Subscription{
// 		UserID:         uint(req.UserId),
// 		SubscriptionID: order["id"].(string),
// 		Active:         "InActive",
// 		StartDate:      startDate,
// 		EndDate:        startDate.AddDate(0, 0, 30),
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = s.Repo.CreatePayment(model.Payment{
// 		SubscriptionID: order["id"].(string),
// 		UserID:         uint(req.UserId),
// 		Amount:         plan.Price,
// 		Status:         "faild",
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &proto.CreateSubscriptionResponse{
// 		Message: "Order id: " + order["id"].(string),
// 		Success: true,
// 	}, nil

// }

// func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, req *proto.UpdatePaymentReq) (*proto.PaymentCommonRes, error) {
// 	err := payment.RazorPaymentVerification(req.Signature, req.OrderId, req.PaymentId)
// 	if err != nil {
// 		log.Println("Faild to update paymetn stauts: ", err)
// 		return nil, err
// 	}
// 	err = s.Repo.UpdateStatus(req.OrderId, req.PaymentId, "success")
// 	if err != nil {
// 		log.Println("Faild to update paymetn stauts: ", err)
// 		return nil, err
// 	}
// 	return &proto.PaymentCommonRes{
// 		Message: "Payment Update Successfully",
// 		Status:  200,
// 	}, nil
// }

// func (s *PaymentService) RenewSubscription(ctx context.Context, req *proto.RenewSubReq) (*proto.PaymentCommonRes, error) {
// 	var wg sync.WaitGroup
// 	var receipt string

// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		receipt = helper.RandString()
// 	}()
// 	fmt.Println("user,id", req.UserId)
// 	userSub, err := s.Repo.GetActiveSubscription(uint(req.UserId))
// 	if err != nil {
// 		return nil, err
// 	}
// 	if userSub.ID == 0 {
// 		return nil, errors.New("no active subscription found")
// 	}
// 	if userSub.EndDate.After(time.Now()) {
// 		return nil, errors.New("current subscription is still active and doesn't need renewal yet")
// 	}
// 	plan, err := s.Repo.GetPlanByID(uint(req.PlanId))
// 	if err != nil {
// 		return nil, err
// 	}
// 	wg.Wait()

// 	orderData := map[string]interface{}{
// 		"amount":   plan.Price * 100,
// 		"currency": "INR",
// 		"receipt":  receipt,
// 	}
// 	order, err := s.RazorClient.Order.Create(orderData, nil)
// 	if err != nil {
// 		log.Println("Faild to create order : ", err)
// 		return nil, err
// 	}
// 	startDate := time.Now()
// 	userSub.OrderID = receipt
// 	userSub.SubscriptionID = order["id"].(string)
// 	userSub.StartDate = startDate
// 	userSub.EndDate = startDate.AddDate(0, 0, 30)
// 	userSub.OrderID = order["id"].(string)

// 	err = s.Repo.UpdateSubscription(userSub)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = s.Repo.CreatePayment(model.Payment{
// 		SubscriptionID: order["id"].(string),
// 		UserID:         uint(req.UserId),
// 		Amount:         plan.Price,
// 		Status:         "faild",
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &proto.PaymentCommonRes{
// 		Message: "Order id: " + order["id"].(string),
// 		Status:  200,
// 	}, nil
// }

func (s *PaymentService) GetSubDetails(ctx context.Context, req *proto.GetSubReq) (*proto.GetSubRes, error) {
	endDate, err := s.Repo.GetSubDetails(uint(req.UserId))
	if err != nil {
		log.Println("Failed to Find SubDetails: ", err.Error())
		return nil, err
	}
	fmt.Println("end = ", endDate.Unix(), " now = ", time.Now().Unix())
	return &proto.GetSubRes{
		ExpiryTime: endDate.Unix(),
	}, nil
}
