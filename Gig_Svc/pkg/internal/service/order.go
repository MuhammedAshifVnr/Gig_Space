package service

import (
	"context"
	"fmt"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/utils/helper"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

func (s *GigService) CreateOrder(ctx context.Context, req *proto.CreateOrderReq) (*proto.CommonGigRes, error) {
	GigRes, err := s.repos.GetGigByID(uint(req.GigId))
	if err != nil {
		log.Println("Failded to find the Gig: ", err.Error())
		return nil, err
	}
	rand := fmt.Sprintf("Odr_%s", helper.RandString())
	Order := model.Order{
		OrderID:      rand,
		GigID:        GigRes.ID,
		ClinetID:     uint(req.ClinetId),
		FreelancerID: GigRes.FreelancerID,
		Amount:       int(GigRes.Price),
	}
	err = s.repos.CreateOrder(Order)
	if err != nil {
		log.Println("Failed to create order: ", err.Error())
		return nil, err
	}
	PaymentRes, err := s.paymetnClient.CreatePaymentOrder(context.Background(), &proto.CreatePaymentOrderReq{
		OrderId: rand,
		Amount:  int64(GigRes.Price),
	})
	if err != nil {
		log.Println("Failed to create payment: ", err.Error())
		return nil, err
	}
	return &proto.CommonGigRes{
		Status:  200,
		Message: "Payment ID: " + PaymentRes.Message,
	}, nil
}
