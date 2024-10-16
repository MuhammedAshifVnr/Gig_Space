package service

import (
	"context"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

func (s *UserService) GetCategoryByName(ctx context.Context, req *proto.CategoryName) (*proto.CategoryIdRes, error) {
	id, err := s.reops.GetCategoryID(req.Name)
	if err != nil {
		log.Println("Failed to Get Category: ",err.Error())
		return &proto.CategoryIdRes{}, err
	}
	return &proto.CategoryIdRes{Id: uint32(id)}, nil
}
