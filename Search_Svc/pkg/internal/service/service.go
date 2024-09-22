package service

import (
	"context"
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Search_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

type SearchService struct {
	repo repo.RepoInter
	proto.UnimplementedSearchServiceServer
}

func NewSearchService(repo repo.RepoInter) *SearchService {
	return &SearchService{
		repo: repo,
	}
}

func (s *SearchService) IndexGig(ctx context.Context, req *proto.IndexGigRequest) (*proto.SearchEmptyRes, error) {
	gig := model.Gig{
		Id:           req.Id,
		Title:        req.Title,
		Description:  req.Description,
		Category:     req.Category,
		Price:        float64(req.Price),
		DeliveryDays: req.DeliveryDays,
		Revisions:    req.Revisions,
		FreelancerId: req.FreelancerId,
		Images:       req.Image,
	}
	fmt.Println("gig id ", req.Id)
	err := s.repo.IndexGig(gig, "gig")
	if err != nil {
		return nil, err
	}
	return &proto.SearchEmptyRes{}, nil
}
