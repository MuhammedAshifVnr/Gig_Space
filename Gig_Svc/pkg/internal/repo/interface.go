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
}
