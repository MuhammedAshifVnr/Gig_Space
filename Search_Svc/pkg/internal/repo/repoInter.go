package repo

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/model"
)

type RepoInter interface {
	IndexGig(data model.Gig, index string) error
	SearchGigs(ctx context.Context, query string, priceUpto float32, revisionsMin, deliveryDaysMax int32) ([]*model.Gig, error)
}
