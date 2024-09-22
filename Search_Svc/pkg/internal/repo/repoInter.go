package repo

import "github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/model"

type RepoInter interface {
	IndexGig(data model.Gig, index string) error
}
