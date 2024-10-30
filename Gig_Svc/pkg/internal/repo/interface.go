package repo

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

type RepoInter interface {
	CreateGgi(gig model.Gig) (model.Gig, error)
	//AddImages(images []string, id uint) error
	GetGigsByFreelancerID(freelancerID uint) ([]model.Gig, error)
	GetGigByID(Id uint) (model.Gig, error)
	UpdateGig(gig model.Gig) error
	DeleteImages(id uint) error
	DeleteGig(id, user_id uint) error
	CreateOrder(data model.Order) error
	GetOrders(clientID uint) ([]*proto.Order, error)
	CreateQuote(Quote model.Quote) error
	GetAllQuotes(freelancerID uint)([]*proto.Quote,error)
	CreateCustomGig(gig model.CustomGig)error
	GetAllOffers(clientID uint) ([]*proto.CreateOfferReq, error) 
	GetCustomGig(ID uint) (model.CustomGig, error)
	CreateCustomOrder(order model.CustomOrder) error
	UpdateOrderStatus(order_id, status string) error
	GetRequest(userid uint) ([]*proto.Request, []*proto.Request, error) 
	UpdateOfferOrderStatus(order_id, status string) error
	AcceptCustomOrder(orderID string)error
	AcceptOrder(orderID string)error
	RejectCustomOrder(orderID string)error
	RejectOrder(orderID string)error
	AdminGetOrders() ([]*proto.Request, []*proto.Request, error) 
	GetCustomOrderByID(order_id string)(model.CustomOrder,error)
	GetOrderByID(order_id string)(model.Order,error)
}
