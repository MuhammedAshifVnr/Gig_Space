package service

import (
	"context"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/internal/model"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/utils/convert"
	"github.com/MuhammedAshifVnr/Gig_Space/Gig_Svc/utils/upimage"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/aws/aws-sdk-go/service/s3"
)

type GigService struct {
	repos      repo.RepoInter
	s3         *s3.S3
	userClient proto.UserServiceClient
	proto.UnimplementedGigServiceServer
}

func NewGigService(repo repo.RepoInter, s3Svc *s3.S3, UserClient proto.UserServiceClient) *GigService {
	return &GigService{
		repos:      repo,
		s3:         s3Svc,
		userClient: UserClient,
	}
}

func (s *GigService) CreateGig(ctx context.Context, req *proto.CreateGigReq) (*proto.EmptyResponse, error) {

	imageUrls := []string{}

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

		imageUrls = append(imageUrls, imageUrl)
	}
	res, err := s.userClient.GetCategoryByName(context.Background(), &proto.CategoryName{Name: req.Category})
	if err != nil {
		return &proto.EmptyResponse{}, err
	}
	gig := model.Gig{
		Title:        req.Title,
		Description:  req.Description,
		Price:        req.Price,
		Category:     uint(res.Id),
		FrelancerID:  uint(req.UserId),
		Image_Urls:   imageUrls,
		DeliveryDays: int(req.DeliveryDays),
		Revisions:    int(req.NumberOfRevisions),
	}
	err = s.repos.CreateGgi(gig)
	if err != nil {
		return &proto.EmptyResponse{}, err
	}
	return &proto.EmptyResponse{}, nil
}

func (s *GigService)GetAllGigByID(ctx context.Context,req *proto.GetAllGigsByIDReq)(*proto.GetAllGigsResp,error){
	gigs,err:=s.repos.GetAllGigByID(uint(req.Id))
	if err!=nil{
		return &proto.GetAllGigsResp{},err
	}
	return &proto.GetAllGigsResp{
		Gig: gigs,
	},nil
}