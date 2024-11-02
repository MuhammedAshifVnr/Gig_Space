package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type RefundEvent struct {
	UserID  uint `json:"user_id"`
	OrderID string
	Event   string
	Amoutn  int
}

func (s *GigService) AdminOrderController(ctx context.Context, req *proto.EmptyGigReq) (*proto.AdOrderController, error) {
	order, Corder, err := s.repos.AdminGetOrders("Freelancer Rejected")
	if err != nil {
		return nil, err
	}
	return &proto.AdOrderController{
		Gigs:      order,
		OfferGigs: Corder,
	}, nil
}

func (s *GigService) AdOrderRefund(ctx context.Context, req *proto.AdRefundReq) (*proto.CommonGigRes, error) {
	var User uint
	var Amount int
	if req.OrderId[0] == 'C' {
		order, err := s.repos.GetCustomOrderByID(req.OrderId)
		if err != nil {
			return nil, err
		}
		_, err = s.paymetnClient.AddWalletRefund(context.Background(), &proto.AddRefund{
			UserId: uint64(order.ClinetID),
			Amount: int64(order.Amount),
		})
		if err != nil {
			if err := s.SendNotification(ctx, RefundEvent{
				UserID:  uint(order.ClinetID),
				OrderID: order.OrderID,
				Event:   "Fail",
				Amoutn:  0,
			}, viper.GetString("RefundTopic")); err != nil {
				return nil, err
			}
			return nil, err
		}
		Amount = order.Amount
		User = order.ClinetID
		err = s.repos.UpdateOrderStatus(order.OrderID, "Cancelled")
		if err != nil {
			return nil, err
		}
	} else {
		order, err := s.repos.GetOrderByID(req.OrderId)
		if err != nil {
			return nil, err
		}
		_, err = s.paymetnClient.AddWalletRefund(context.Background(), &proto.AddRefund{
			UserId: uint64(order.ClinetID),
			Amount: int64(order.Amount),
		})
		if err != nil {
			if err := s.SendNotification(ctx, RefundEvent{
				UserID:  uint(order.ClinetID),
				OrderID: order.OrderID,
				Event:   "Fail",
				Amoutn:  0,
			}, viper.GetString("RefundTopic")); err != nil {
				return nil, err
			}
			return nil, err
		}
		Amount = order.Amount
		User = order.ClinetID
		err = s.repos.UpdateOrderStatus(order.OrderID, "Cancelled")
		if err != nil {
			return nil, err
		}
	}
	if err := s.SendNotification(ctx, RefundEvent{
		UserID:  uint(User),
		OrderID: req.OrderId,
		Event:   "Done",
		Amoutn:  Amount,
	}, viper.GetString("RefundTopic")); err != nil {
		return nil, err
	}
	return &proto.CommonGigRes{
		Message: "Refund Done",
		Status:  200,
	}, nil
}

func (s *GigService) SendNotification(ctx context.Context, event RefundEvent, key string) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	writer, ok := s.kafkaWriter[key]
	if ok {
		err = writer.WriteMessages(ctx,
			kafka.Message{
				Key:   []byte(key),
				Value: eventData,
			},
		)
	}
	if err != nil {
		return fmt.Errorf("failed to publish notification: %w", err)
	}
	return nil
}

func (s *GigService) AdOrderCheck(ctx context.Context, req *proto.EmptyGigReq) (*proto.AdOrderController, error) {
	order, Corder, err := s.repos.AdminGetOrders("Done")
	if err != nil {
		return nil, err
	}
	return &proto.AdOrderController{
		Gigs:      order,
		OfferGigs: Corder,
	}, nil
}

func(s *GigService)AdPaymentTransfer(ctx context.Context,req *proto.PaymentTransferReq)(*proto.CommonGigRes,error){
	var User uint
	var Amount int
	if req.OrderId[0] == 'C' {
		order, err := s.repos.GetCustomOrderByID(req.OrderId)
		if err != nil {
			return nil, err
		}
		_, err = s.paymetnClient.AddWalletRefund(context.Background(), &proto.AddRefund{
			UserId: uint64(order.ClinetID),
			Amount: int64(order.Amount),
		})
		if err != nil {
			if err := s.SendNotification(ctx, RefundEvent{
				UserID:  uint(order.ClinetID),
				OrderID: order.OrderID,
				Event:   "Fail",
				Amoutn:  0,
			}, viper.GetString("RefundTopic")); err != nil {
				return nil, err
			}
			return nil, err
		}
		Amount = order.Amount
		User = order.ClinetID
		err = s.repos.UpdateOrderStatus(order.OrderID, "Order Completed Successfully")
		if err != nil {
			return nil, err
		}
	} else {
		order, err := s.repos.GetOrderByID(req.OrderId)
		if err != nil {
			return nil, err
		}
		_, err = s.paymetnClient.AddWalletRefund(context.Background(), &proto.AddRefund{
			UserId: uint64(order.ClinetID),
			Amount: int64(order.Amount),
		})
		if err != nil {
			if err := s.SendNotification(ctx, RefundEvent{
				UserID:  uint(order.ClinetID),
				OrderID: order.OrderID,
				Event:   "Fail",
				Amoutn:  0,
			}, viper.GetString("RefundTopic")); err != nil {
				return nil, err
			}
			return nil, err
		}
		Amount = order.Amount
		User = order.ClinetID
		err = s.repos.UpdateOrderStatus(order.OrderID, "Order Completed Successfully")
		if err != nil {
			return nil, err
		}
	}
	if err := s.SendNotification(ctx, RefundEvent{
		UserID:  uint(User),
		OrderID: req.OrderId,
		Event:   "Done",
		Amoutn:  Amount,
	}, viper.GetString("RefundTopic")); err != nil {
		return nil, err
	}
	return &proto.CommonGigRes{
		Message: "Payment Transfer Done",
		Status:  200,
	}, nil
}