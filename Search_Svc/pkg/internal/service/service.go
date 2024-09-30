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

func (s *SearchService) SearchGig(ctx context.Context, req *proto.SearchGigReq) (*proto.SearchGigRes, error) {
	gigs, err := s.repo.SearchGigs(ctx, req.Query, req.PriceUpto, req.RevisionsMin, req.DeliveryDaysMax)
	if err != nil {
		return nil, err
	}
	var responseGigs []*proto.ResultGig
	for _, gig := range gigs {
		responseGigs = append(responseGigs, &proto.ResultGig{
			Id:           uint32(gig.Id),
			Title:        gig.Title,
			Description:  gig.Description,
			Image:        gig.Images,
			Price:        float32(gig.Price),
			Revisions:    gig.Revisions,
			DeliveryDays: gig.DeliveryDays,
		})
	}

	return &proto.SearchGigRes{Gigs: responseGigs}, nil
}

func (s *SearchService) UpdateIndexGig(ctx context.Context, req *proto.IndexGigRequest) (*proto.SearchEmptyRes, error) {
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
	err := s.repo.UpdateGig(gig, "gig")
	if err != nil {
		return nil, err
	}
	return &proto.SearchEmptyRes{}, nil
}

func (s *SearchService) DeleteGig(ctx context.Context, req *proto.DeleteDocumentReq) (*proto.SearchEmptyRes,error) {
	docID, err := s.repo.GetDocId(uint(req.Id), "gig")
	if err != nil {
		return nil, err
	}
	err = s.repo.DeleteDocument(docID, "gig")
	if err != nil {
		return nil, err
	}
	return &proto.SearchEmptyRes{}, nil
}
