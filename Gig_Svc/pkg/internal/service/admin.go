package service

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

func (s *GigService) AdminOrderController(ctx context.Context, req *proto.EmptyReq) (*proto.AdOrderController, error) {
	order, Corder, err := s.repos.AdminGetOrders()
	if err != nil {
		return nil, err
	}
	return &proto.AdOrderController{
		Gigs:      order,
		OfferGigs: Corder,
	}, nil
}

func (s *GigService) AdOrderRefund(ctx context.Context, req *proto.AdRefundReq) (*proto.CommonGigRes, error) {
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
			
			return nil, err
		}
	}
	return &proto.CommonGigRes{},nil
}
