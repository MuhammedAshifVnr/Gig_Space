package service

import (
	"context"
	"fmt"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/utils/helper"
	"github.com/spf13/viper"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

type StatusEvent struct {
	OrderID string
	Event   string
}

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
		Status:       "Payment is Pending",
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

func (s *GigService) GetClientOrders(ctx context.Context, req *proto.GetOrderReq) (*proto.GetOrderRes, error) {
	result, err := s.repos.GetOrders(uint(req.UserId))
	if err != nil {
		return nil, err
	}
	return &proto.GetOrderRes{
		Orders: result,
	}, nil
}

func (s *GigService) RequestQuote(ctx context.Context, req *proto.QuoteReq) (*proto.CommonGigRes, error) {
	fmt.Println(req)
	err := s.repos.CreateQuote(model.Quote{
		GigId:        req.GigId,
		ClientId:     req.ClientId,
		Describe:     req.Describe,
		Price:        req.Price,
		DeliveryDays: int(req.DeliveryDays),
	})
	if err != nil {
		return &proto.CommonGigRes{}, err
	}
	return &proto.CommonGigRes{
		Message: "Request sent successfully!",
		Status:  200,
	}, nil
}

func (s *GigService) GetAllQuotes(ctx context.Context, req *proto.GetAllQuoteReq) (*proto.GetAllQuoteRes, error) {
	quotes, err := s.repos.GetAllQuotes(uint(req.UserId))
	if err != nil {
		return nil, err
	}
	fmt.Println(quotes)
	return &proto.GetAllQuoteRes{
		Quotes: quotes,
	}, nil
}

func (s *GigService) CreateOffer(ctx context.Context, req *proto.CreateOfferReq) (*proto.CommonGigRes, error) {
	err := s.repos.CreateCustomGig(model.CustomGig{
		GigRequestID: uint(req.GigRequestId),
		FreelancerID: uint(req.FreelancerId),
		ClientID:     uint(req.ClientId),
		Title:        req.Title,
		Description:  req.Descripton,
		Price:        float64(req.Price),
		DeliveryDays: int(req.DeliveryDays),
	})
	if err != nil {
		log.Println("Faild to create the CustomGig: ", err.Error())
		return nil, err
	}
	return &proto.CommonGigRes{
		Message: "CustomGig sent successfully!",
		Status:  200,
	}, nil
}

func (s *GigService) GetAllOffers(ctx context.Context, req *proto.GetAllOfferReq) (*proto.GetAllOfferRes, error) {
	offers, err := s.repos.GetAllOffers(uint(req.ClientId))
	if err != nil {
		log.Println("Faild to find the offers: ", err.Error())
		return nil, err
	}
	return &proto.GetAllOfferRes{
		Offers: offers,
	}, nil
}

func (s *GigService) CreateOfferOrder(ctx context.Context, req *proto.CreateOrderReq) (*proto.CommonGigRes, error) {
	Gig, err := s.repos.GetCustomGig(uint(req.GigId))
	if err != nil {
		log.Println("Faild to find the offer: ", err.Error())
		return nil, err
	}
	rand := fmt.Sprintf("Codr_%s", helper.RandString())
	err = s.repos.CreateCustomOrder(model.CustomOrder{
		OrderID:      rand,
		CustomGigID:  Gig.ID,
		ClinetID:     uint(req.ClinetId),
		FreelancerID: Gig.FreelancerID,
		Status:       "Payment is Pending",
		Amount:       int(Gig.Price),
	})
	if err != nil {
		log.Println("Faild to Create the order: ", err.Error())
		return nil, err
	}
	PaymentRes, err := s.paymetnClient.CreatePaymentOrder(context.Background(), &proto.CreatePaymentOrderReq{
		OrderId: rand,
		Amount:  int64(Gig.Price),
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

func (s *GigService) UpdateOrderStatus(ctx context.Context, req *proto.OrderStatusReq) (*proto.CommonGigRes, error) {
	fmt.Println(req.OrderId[0])
	if req.OrderId[0] == 'C' {
		err := s.repos.UpdateOfferOrderStatus(req.OrderId, req.Status)
		if err != nil {
			return nil, err
		}
	} else {
		err := s.repos.UpdateOrderStatus(req.OrderId, req.Status)
		if err != nil {
			return nil, err
		}
	}
	return &proto.CommonGigRes{
		Message: "Status updated",
		Status:  200,
	}, nil
}

func (s *GigService) GetAllRequest(ctx context.Context, req *proto.GetAllRequestReq) (*proto.GetAllRequestRes, error) {
	orders, custom_orders, err := s.repos.GetRequest(uint(req.UserId))
	if err != nil {
		return nil, err
	}
	fmt.Println(req.UserId)
	return &proto.GetAllRequestRes{
		Gigs:      orders,
		OfferGigs: custom_orders,
	}, nil
}

func (s *GigService) AcceptRequest(ctx context.Context, req *proto.AcceptReq) (*proto.CommonGigRes, error) {
	var err error
	if req.OrderId[0] == 'C' {
		err = s.repos.AcceptCustomOrder(req.OrderId)
	} else {
		err = s.repos.AcceptOrder(req.OrderId)
	}

	if err != nil {
		log.Printf("Failed to accept order %s: %v", req.OrderId, err)
		return nil, err
	}

	if err := s.SendNotification(ctx, model.OrderEvent{
		OrderID: req.OrderId,
		Event:   "OrderAccepted",
	}, viper.GetString("OrderTopic")); err != nil {
		return nil, err
	}

	log.Printf("Order %s accepted successfully", req.OrderId)
	return &proto.CommonGigRes{
		Message: "Order Accepted.",
		Status:  200,
	}, nil
}

func (s *GigService) RejectRequest(ctx context.Context, req *proto.RejectReq) (*proto.CommonGigRes, error) {
	var err error
	if req.OrderId[0] == 'C' {
		err = s.repos.RejectCustomOrder(req.OrderId)
	} else {
		err = s.repos.RejectOrder(req.OrderId)
	}

	if err != nil {
		log.Printf("Failed to reject order %s: %v", req.OrderId, err)
		return nil, err
	}

	if err := s.SendNotification(ctx, model.OrderEvent{
		OrderID: req.OrderId,
		Event:   "OrderRejection",
	}, viper.GetString("OrderTopic")); err != nil {
		return nil, err
	}

	log.Printf("Order %s rejected successfully", req.OrderId)
	return &proto.CommonGigRes{
		Message: "Order Rejected.",
		Status:  200,
	}, nil
}

func (s *GigService) GetAllOrders(ctx context.Context, req *proto.AllOrdersReq) (*proto.AllOrdersRes, error) {
	order, err := s.repos.GetAllOrders(uint(req.UserId))
	if err != nil {
		log.Printf("Failed to find order: %v", err)
		return nil, err
	}
	COrder, err := s.repos.GetAllCustomOrders(uint(req.UserId))
	if err != nil {
		log.Printf("Failed to find order: %v", err)
		return nil, err
	}
	return &proto.AllOrdersRes{
		Orders:  order,
		COrders: COrder,
	}, nil
}

func (s *GigService) GetOrderByID(ctx context.Context, req *proto.OrderByIDReq) (*proto.OrderDetail, error) {
	if req.OrderId[0] == 'C' {
		order, err := s.repos.GetCustomOrderDetail(req.OrderId)
		if err != nil {
			log.Printf("Failed to find order: %v", err)
			return nil, err
		}
		return &proto.OrderDetail{
			OrderId:      order.OrderID,
			GigId:        uint64(order.CustomGigID),
			Status:       order.Status,
			ClientId:     uint64(order.ClinetID),
			FrelancerId:  uint64(order.FreelancerID),
			Amount:       int64(order.Amount),
			LastUpdated:  order.UpdatedAt.String(),
			OrderCreated: order.CreatedAt.String(),
		}, nil
	} else {
		order, err := s.repos.GetOrderDetail(req.OrderId)
		if err != nil {
			log.Printf("Failed to find order: %v", err)
			return nil, err
		}
		return &proto.OrderDetail{
			OrderId:      order.OrderID,
			GigId:        uint64(order.GigID),
			Status:       order.Status,
			ClientId:     uint64(order.ClinetID),
			FrelancerId:  uint64(order.FreelancerID),
			Amount:       int64(order.Amount),
			LastUpdated:  order.UpdatedAt.String(),
			OrderCreated: order.CreatedAt.String(),
		}, nil
	}
}

func (s *GigService) ClientUpdatePendingStatus(ctx context.Context, req *proto.OrderIDReq) (*proto.CommonGigRes, error) {
	if req.OrderId[0] == 'C' {
		err := s.repos.CordrUpdatePendingStatus(req.OrderId, uint(req.ClientId))
		if err != nil {
			return nil, err
		}
	} else {
		err := s.repos.OrderUpdatePendingStatus(req.OrderId, uint(req.ClientId))
		if err != nil {
			return nil, err
		}
	}
	if err := s.SendNotification(ctx, StatusEvent{
		OrderID: req.OrderId,
		Event:   "Pending",
	}, viper.GetString("StatusTopic")); err != nil {
		return nil, err
	}
	return &proto.CommonGigRes{
		Message: "Order Updated",
		Status:  200,
	}, nil
}

func (s *GigService) ClientUpdateDoneStatus(ctx context.Context, req *proto.OrderIDReq) (*proto.CommonGigRes, error) {
	if req.OrderId[0] == 'C' {
		err := s.repos.CordrUpdateDoneStatus(req.OrderId, uint(req.ClientId))
		if err != nil {
			return nil, err
		}
	} else {
		err := s.repos.OrderUpdateDoneStatus(req.OrderId, uint(req.ClientId))
		if err != nil {
			return nil, err
		}
	}
	if err := s.SendNotification(ctx, RefundEvent{
		OrderID: req.OrderId,
		Event:   "Done",
	}, viper.GetString("StatusTopic")); err != nil {
		return nil, err
	}
	return &proto.CommonGigRes{
		Message: "Order Updated",
		Status:  200,
	}, nil
}

func (s *GigService) GetFreelancerIDByOrder(ctx context.Context, req *proto.OrderByIDReq) (*proto.UserIDRes, error) {
	if req.OrderId[0] == 'C' {
		order, err := s.repos.GetCustomOrderDetail(req.OrderId)
		if err != nil {
			return nil, err
		}
		return &proto.UserIDRes{
			UserId: uint64(order.FreelancerID),
		}, nil
	} else {
		order, err := s.repos.GetOrderByID(req.OrderId)
		if err != nil {
			return nil, err
		}
		return &proto.UserIDRes{
			UserId: uint64(order.FreelancerID),
		}, nil
	}
}
