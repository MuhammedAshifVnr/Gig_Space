package repo

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/model"
)

type RepoInter interface {
	IndexGig(data model.Gig, index string) error
	UpdateGig(data model.Gig,index string) error
	DeleteDocument(docID, index string) error 
	GetDocId(gigID uint, index string) (string, error)
	SearchGigs(ctx context.Context, query string, priceUpto float32, revisionsMin, deliveryDaysMax int32) ([]*model.Gig, error)
}
