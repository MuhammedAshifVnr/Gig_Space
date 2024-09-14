package repo

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/internal/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

type RepoInter interface {
	CreateGgi(gig model.Gig) error
	GetAllGigByID(id uint) ([]*proto.Gig, error)
}
