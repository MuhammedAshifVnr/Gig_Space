package service

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

func (s *UserService) UpdateBio(ctc context.Context, req *proto.UpdateProfileReq) (*proto.CommonRes, error) {
	user, _ := s.reops.GetBio(uint(req.UserId))
	if user.ID == 0 {
		err := s.reops.CreateBio(model.UserProfile{
			UserID: uint(req.UserId),
			Bio:    req.Description,
			Title:  req.Title,
		})
		if err != nil {
			return &proto.CommonRes{}, err
		}
	} else {
		err := s.reops.UpdateBio(model.UserProfile{
			UserID: uint(req.UserId),
			Bio:    req.Description,
			Title:  req.Title,
		})
		if err != nil {
			return &proto.CommonRes{}, err
		}
	}
	return &proto.CommonRes{
		Message: "Profile Updated.",
		Status:  200,
	}, nil
}

func (s *UserService) FreelacerAddSkill(ctx context.Context, req *proto.FreeAddSkillsReq) (*proto.CommonRes, error) {

	val, err := s.reops.GetSkill(req.SkillName)
	if err != nil {
		return &proto.CommonRes{}, err
	}
	if val.ID == 0 {
		return &proto.CommonRes{
			Message: "Please enter valid skill.",
			Status:  400,
		}, nil
	}
	err = s.reops.FreelacerAddSkill(uint(req.UserId), val.ID, int(req.ProficiencyLevel))
	if err != nil {
		return &proto.CommonRes{}, err
	}

	return &proto.CommonRes{
		Message: "Skills Updated.",
		Status:  200,
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *proto.ProfileReq) (*proto.ProfileRes, error) {
	user, err := s.reops.GetProfile(uint(req.UserId))
	if err != nil {
		return &proto.ProfileRes{}, err
	}
	skills, err := s.reops.GetSkillByuserID(uint(req.UserId))
	if err != nil {
		return &proto.ProfileRes{}, err
	}
	return &proto.ProfileRes{
		Firstname:   user.FirstName,
		Lastname:    user.LastName,
		Email:       user.Email,
		Country:     user.Country,
		Phone:       user.Phone,
		Skill:       skills,
		Title:       user.Title,
		Description: user.Bio,
	}, nil
}

func (s *UserService) DeleteSkill(ctx context.Context,req *proto.DeleteSkillRes) (*proto.CommonRes,error){
	err:=s.reops.DeleteSkillByID(uint(req.UserId),uint(req.SkillId))
	if err != nil{
		return &proto.CommonRes{},err
	}
	return &proto.CommonRes{
		Message: "Skill Deleted Successfully",
		Status: 200,
	},nil
}