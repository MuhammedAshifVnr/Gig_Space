package service

import (
	"context"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/utils/convert"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/utils/upimage"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/aws/aws-sdk-go/service/s3"
)

type GigService struct {
	repos         repo.RepoInter
	s3            *s3.S3
	userClient    proto.UserServiceClient
	searchClient  proto.SearchServiceClient
	paymetnClient proto.PaymentServiceClient
	proto.UnimplementedGigServiceServer
}

func NewGigService(repo repo.RepoInter, s3Svc *s3.S3, UserClient proto.UserServiceClient, SearchClient proto.SearchServiceClient,payment proto.PaymentServiceClient) *GigService {
	return &GigService{
		repos:        repo,
		s3:           s3Svc,
		searchClient: SearchClient,
		userClient:   UserClient,
		paymetnClient: payment,
	}
}

func (s *GigService) CreateGig(ctx context.Context, req *proto.CreateGigReq) (*proto.EmptyResponse, error) {
	_, err := s.userClient.GetCategoryByName(context.Background(), &proto.CategoryName{Name: req.Category})
	if err != nil {
		log.Println("Faild to Find Category: ", err.Error())
		return &proto.EmptyResponse{}, err
	}
	gig := model.Gig{
		Title:        req.Title,
		Description:  req.Description,
		FreelancerID: uint(req.UserId),
		Category:     req.Category,
		Price:        req.Price,
		DeliveryDays: int(req.DeliveryDays),
		Revisions:    int(req.NumberOfRevisions),
	}
	for _, imageBytes := range req.GetImages() {
		file, fileHeader, err := convert.ConvertToMultipartFile(imageBytes)
		if err != nil {
			log.Println("Error converting image to multipart:", err)
			return nil, err
		}

		imageUrl, err := upimage.UploadImage(s.s3, file, fileHeader)
		if err != nil {
			log.Println("Error uploading image to S3:", err)
			return nil, err
		}

		gig.Images = append(gig.Images, model.Image{Url: imageUrl})
	}
	resGig, err := s.repos.CreateGgi(gig)
	if err != nil {
		log.Println("Failed to Create Gig: ", err.Error())
		return &proto.EmptyResponse{}, err
	}
	_, err = s.searchClient.IndexGig(context.Background(), &proto.IndexGigRequest{
		Id:           uint64(resGig.ID),
		Title:        gig.Title,
		Description:  gig.Description,
		Category:     req.Category,
		Price:        float32(gig.Price),
		DeliveryDays: int32(gig.DeliveryDays),
		Revisions:    int32(gig.Revisions),
		FreelancerId: uint64(gig.FreelancerID),
		Image:        gig.Images[0].Url,
	})
	if err != nil {
		delErr := s.repos.DeleteGig(resGig.ID, gig.FreelancerID)
		if delErr != nil {
			log.Println("Failed to rollback database entry after Elasticsearch failure: ", delErr.Error())
			return &proto.EmptyResponse{}, delErr
		}
		log.Println("Failed to Create document in elastic : ", err.Error())
		return nil, err
	}
	return &proto.EmptyResponse{}, nil
}

func (s *GigService) GetGigsByFreelancerID(ctx context.Context, req *proto.GetGigsByFreelancerIDRequest) (*proto.GetGigsByFreelancerIDResponse, error) {
	gigs, err := s.repos.GetGigsByFreelancerID(uint(req.FreelancerId))
	if err != nil {
		log.Println("Failed to Find User Gigs: ", err.Error())
		return &proto.GetGigsByFreelancerIDResponse{}, err
	}
	var grpcGigs []*proto.Gig
	for _, gig := range gigs {
		var grpcImages []*proto.Image
		for _, img := range gig.Images {
			grpcImages = append(grpcImages, &proto.Image{
				Url: img.Url,
			})
		}
		grpcGigs = append(grpcGigs, &proto.Gig{
			Id:           uint64(gig.ID),
			Title:        gig.Title,
			Description:  gig.Description,
			Category:     gig.Category,
			FreelancerId: uint64(gig.FreelancerID),
			Price:        float32(gig.Price),
			DeliveryDays: int32(gig.DeliveryDays),
			Revisions:    int32(gig.Revisions),
			Image:        grpcImages,
		})
	}

	return &proto.GetGigsByFreelancerIDResponse{Gigs: grpcGigs}, nil
}

func (s *GigService) UpdateGigByID(ctx context.Context, req *proto.UpdateGigRequest) (*proto.CommonGigRes, error) {
	gig, err := s.repos.GetGigByID(uint(req.Id))
	if err != nil {
		log.Println("Failed to Find Gig: ", err.Error())
		return nil, err
	}
	if req.Category != "" {
		_, err := s.userClient.GetCategoryByName(context.Background(), &proto.CategoryName{Name: req.Category})
		if err != nil {
			log.Println("Failed to Find Category: ", err.Error())
			return nil, err
		}
		gig.Category = req.Category
	}

	if req.Title != "" {
		gig.Title = req.Title
	}
	if req.Description != "" {
		gig.Description = req.Description
	}
	if req.Price != 0 {
		gig.Price = req.Price
	}
	if req.DeliveryDays != 0 {
		gig.DeliveryDays = int(req.DeliveryDays)
	}
	if req.NumberOfRevisions != 0 {
		gig.Revisions = int(req.NumberOfRevisions)
	}
	if len(req.Images) > 0 {
		if err = s.repos.DeleteImages(gig.ID); err != nil {
			return nil, err
		}
		gig.Images = []model.Image{}
		for _, imageBytes := range req.GetImages() {
			file, fileHeader, err := convert.ConvertToMultipartFile(imageBytes)
			if err != nil {
				return nil, err
			}
			imageUrl, err := upimage.UploadImage(s.s3, file, fileHeader)
			if err != nil {
				return nil, err
			}
			gig.Images = append(gig.Images, model.Image{Url: imageUrl})
		}
	}

	err = s.repos.UpdateGig(gig)
	if err != nil {
		log.Println("Failed to Update Gig: ", err.Error())
		return nil, err
	}
	_, err = s.searchClient.UpdateIndexGig(context.Background(), &proto.IndexGigRequest{
		Id:           uint64(gig.ID),
		Title:        gig.Title,
		Description:  gig.Description,
		Category:     gig.Category,
		Price:        float32(gig.Price),
		DeliveryDays: int32(gig.DeliveryDays),
		Revisions:    int32(gig.Revisions),
		FreelancerId: uint64(gig.FreelancerID),
		Image:        gig.Images[0].Url,
	})
	if err != nil {
		log.Println("Failed to Update ES Document: ", err.Error())
		return nil, err
	}
	return &proto.CommonGigRes{
		Message: "Updated Successfully",
		Status:  200,
	}, nil

}

func (s *GigService) DeleteGigByID(ctx context.Context, req *proto.DeleteReq) (*proto.CommonGigRes, error) {
	err := s.repos.DeleteGig(uint(req.GigId), uint(req.UserId))
	if err != nil {
		log.Println("Failed to delete Gig: ", err.Error())
		return nil, err
	}

	_, err = s.searchClient.DeleteGig(context.Background(), &proto.DeleteDocumentReq{
		Id: req.GigId,
	})
	if err != nil {
		log.Println("Failed to Delete ES documetn: ", err.Error())
		return nil, err
	}
	return &proto.CommonGigRes{
		Message: "Successfully Deleted",
		Status:  200,
	}, nil
}

