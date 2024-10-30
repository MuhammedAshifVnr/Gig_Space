package service

import (
	"context"
	"fmt"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/Payment_Svc/utils/payment"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

func (s *PaymentService) CreatePaymentOrder(ctx context.Context, req *proto.CreatePaymentOrderReq) (*proto.PaymentCommonRes, error) {
	paymentOrder, err := s.RazorClient.Order.Create(map[string]interface{}{
		"amount":   req.Amount * 100,
		"currency": "INR",
		"receipt":  req.OrderId,
	}, nil)
	if err != nil {
		log.Println("Failed to create payment : ", err)
		return nil, err
	}
	payment := model.OrderPayment{
		ReceiptID: req.OrderId,
		OrderID:   paymentOrder["id"].(string),
		Status:    "Pending",
		Amount:    int(req.Amount),
	}
	err = s.Repo.CreateOrderPayment(payment)
	if err != nil {
		log.Println("Failed to Save paymet: ", err.Error())
		return nil, err
	}
	return &proto.PaymentCommonRes{
		Message: paymentOrder["id"].(string),
	}, nil
}

func (s *PaymentService) UpdatePaymentStatus(ctx context.Context, req *proto.UpdatePaymentReq) (*proto.PaymentCommonRes, error) {
	err := payment.RazorPaymentVerification(req.Signature, req.OrderId, req.PaymentId)
	if err != nil {
		log.Println("Failed to verify payment signature:", err)
		return nil, fmt.Errorf("payment verification failed: %w", err)
	}

	receiptID, err := s.Repo.UpdateStatus(req.OrderId, req.PaymentId, "success")
	if err != nil {
		log.Println("Failed to update payment status:", err)
		return nil, fmt.Errorf("failed to update payment status: %w", err)
	}

	_, err = s.GigClient.UpdateOrderStatus(ctx, &proto.OrderStatusReq{
		OrderId: receiptID,
		Status:  "Freelance Approval Pending",
	})
	if err != nil {
		log.Println("Failed to update order status in Gig Service:", err)
		return nil, fmt.Errorf("failed to update order status in gig service: %w", err)
	}

	return &proto.PaymentCommonRes{
		Message: "Payment Update Successfully",
		Status:  200,
	}, nil
}
