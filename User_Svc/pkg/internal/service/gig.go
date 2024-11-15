package service

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/sirupsen/logrus"
)

func (s *UserService) GetCategoryByName(ctx context.Context, req *proto.CategoryName) (*proto.CategoryIdRes, error) {
	id, err := s.reops.GetCategoryID(req.Name)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"category_name": req.Name,
			"error":         err.Error(),
		}).Error("Failed to get category")
		return &proto.CategoryIdRes{}, err
	}
	return &proto.CategoryIdRes{Id: uint32(id)}, nil
}
