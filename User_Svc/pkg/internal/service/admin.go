package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/jwt"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

func (s *UserService) AdminLogin(ctx context.Context, req *proto.AdLoginReq) (*proto.CommonRes, error) {
	admin, err := s.reops.GetAdmin(req.Email)
	if err != nil {
		log.Println("Failed to Find Admin : ", err.Error())
		return &proto.CommonRes{
			Message: "Admin not found",
			Status:  400,
			Error:   err.Error(),
		}, nil
	}
	if admin.Password != req.Password {
		log.Println("Failed to Match Password ")
		return &proto.CommonRes{
			Message: "Password not match.",
			Status:  http.StatusUnauthorized,
		}, nil
	}
	token, err := jwt.GenerateJwtToken(admin.Email, admin.ID, "admin", 0)
	if err != nil {
		log.Println("Failed to Genereate Jwt: ", err.Error())
		return &proto.CommonRes{
			Message: "Error form jwt creation ",
			Status:  404,
			Error:   err.Error(),
		}, nil
	}
	return &proto.CommonRes{
		Message: "Logged in successfully",
		Status:  200,
		Data: map[string]*proto.AnyValue{
			"token": {
				Value: &proto.AnyValue_StringValue{StringValue: token},
			},
		},
	}, nil
}

func (s *UserService) AddCategory(ctx context.Context, req *proto.CategoryReq) (*proto.CommonRes, error) {
	category := model.Category{
		Name:     req.Name,
		IsActive: true,
	}
	err := s.reops.AddCategory(category)
	if err != nil {
		log.Println("Failed to Add Category : ", err.Error())
		return &proto.CommonRes{}, err
	}
	return &proto.CommonRes{
		Message: "Successfully add category",
		Status:  http.StatusAccepted,
	}, nil
}

func (s *UserService) AddSkill(ctx context.Context, req *proto.AddSkillReq) (*proto.CommonRes, error) {
	skill := model.Skills{
		SkillName: req.SkillName,
	}
	err := s.reops.AddSkill(skill)
	if err != nil {
		log.Println("Failed to Add Skill: ", err.Error())
		return &proto.CommonRes{}, err
	}
	return &proto.CommonRes{
		Message: "Skill added Successfully.",
		Status:  200,
	}, nil
}

func (s *UserService) GetCategory(ctx context.Context, req *proto.EmtpyReq) (*proto.GetCategoryRes, error) {
	categors, err := s.reops.GetCategory()
	if err != nil {
		log.Println("Failed to Find Category: ", err.Error())
		return &proto.GetCategoryRes{}, err
	}
	return &proto.GetCategoryRes{
		Category: categors,
	}, nil

}

func (s *UserService) GetSkill(ctx context.Context, req *proto.EmtpyReq) (*proto.GetSkillsRes, error) {
	skills, err := s.reops.GetSkills()
	if err != nil {
		log.Println("Failed to Find Skill: ", err.Error())
		return &proto.GetSkillsRes{}, err
	}
	return &proto.GetSkillsRes{
		Skill: skills,
	}, nil
}

func (s *UserService) AdDeleteSkill(ctx context.Context, req *proto.ADeleteSkillReq) (*proto.EmtpyRes, error) {
	err := s.reops.AdminDeleteSkill(uint(req.Id))
	if err != nil {
		log.Println("Failed to Delete Skill: ", err.Error())
		return &proto.EmtpyRes{}, err
	}
	return &proto.EmtpyRes{}, nil
}

func (s *UserService) DeleteCategory(ctx context.Context, req *proto.DeleteCatReq) (*proto.EmtpyRes, error) {
	err := s.reops.AdminDeleteCategory(uint(req.Id))
	if err != nil {
		log.Println("Failed to Delete Category: ", err.Error())
		return &proto.EmtpyRes{}, err
	}
	return &proto.EmtpyRes{}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *proto.EmtpyReq) (*proto.GetAllUserRes, error) {
	users, err := s.reops.GetAllUsers()
	if err != nil {
		log.Println("Failed to Find Users: ", err.Error())
		return &proto.GetAllUserRes{}, err
	}
	return &proto.GetAllUserRes{
		Users: users,
	}, nil
}

func (s *UserService) AdminUserBlock(ctx context.Context, req *proto.BlockReq) (*proto.CommonRes, error) {
	user, err := s.reops.GetUserByID(uint(req.Id))
	if err != nil {
		log.Println("Failed to Find User: ", err.Error())
		return &proto.CommonRes{}, err
	}
	if !user.IsActive {
		err := s.reops.UnBlockUser(user.ID)
		if err != nil {
			log.Println("Failed to UnBlock User: ", err.Error())
			return &proto.CommonRes{}, err
		}
	} else {
		err := s.reops.BlockUser(user.ID)
		if err != nil {
			log.Println("Failed to Bolck User: ", err.Error())
			return &proto.CommonRes{}, err
		}
	}
	return &proto.CommonRes{
		Message: "Updated Successfully",
		Status:  200,
	}, nil

}

func (s *UserService) GetUserEmail(ctx context.Context, req *proto.ProfileReq) (*proto.EmailRes, error) {
	fmt.Println("USerID",req.UserId)
	User, err := s.reops.GetUserByID(uint(req.UserId))
	if err != nil {
		return nil, err
	}
	fmt.Println(User.Email)
	return &proto.EmailRes{
		Email: User.Email,
	}, nil
}
