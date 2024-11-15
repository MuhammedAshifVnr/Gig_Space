package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/jwt"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/sirupsen/logrus"
)

func (s *UserService) AdminLogin(ctx context.Context, req *proto.AdLoginReq) (*proto.CommonRes, error) {
	admin, err := s.reops.GetAdmin(req.Email)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err.Error(),
		}).Error("Failed to find admin")
		return &proto.CommonRes{
			Message: "Admin not found",
			Status:  400,
			Error:   err.Error(),
		}, nil
	}
	if admin.Password != req.Password {
		s.Log.WithFields(logrus.Fields{
			"email": req.Email,
		}).Warn("Password mismatch")
		return &proto.CommonRes{
			Message: "Password not match.",
			Status:  http.StatusUnauthorized,
		}, nil
	}
	token, err := jwt.GenerateJwtToken(admin.Email, admin.ID, "admin", 0)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err.Error(),
		}).Error("Failed to generate JWT")
		return &proto.CommonRes{
			Message: "Error form jwt creation ",
			Status:  404,
			Error:   err.Error(),
		}, nil
	}

	s.Log.WithFields(logrus.Fields{
		"email":    req.Email,
		"admin_id": admin.ID,
	}).Info("Admin logged in successfully")

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
		s.Log.WithFields(logrus.Fields{
			"category_name": req.Name,
			"error":         err.Error(),
		}).Error("Failed to add category")
		return &proto.CommonRes{}, err
	}

	s.Log.WithField("category_name", req.Name).Info("Category added successfully")
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
		s.Log.WithFields(logrus.Fields{
			"skill_name": req.SkillName,
			"error":      err.Error(),
		}).Error("Failed to add skill")
		return &proto.CommonRes{}, err
	}

	s.Log.WithField("skill_name", req.SkillName).Info("Skill added successfully")
	return &proto.CommonRes{
		Message: "Skill added Successfully.",
		Status:  200,
	}, nil
}

func (s *UserService) GetCategory(ctx context.Context, req *proto.EmtpyReq) (*proto.GetCategoryRes, error) {
	categors, err := s.reops.GetCategory()
	if err != nil {
		s.Log.WithField("error", err.Error()).Error("Failed to find categories")
		return &proto.GetCategoryRes{}, err
	}

	s.Log.Info("Categories retrieved successfully")
	return &proto.GetCategoryRes{
		Category: categors,
	}, nil

}

func (s *UserService) GetSkill(ctx context.Context, req *proto.EmtpyReq) (*proto.GetSkillsRes, error) {
	skills, err := s.reops.GetSkills()
	if err != nil {
		s.Log.WithField("error", err.Error()).Error("Failed to find skills")
		return &proto.GetSkillsRes{}, err
	}

	s.Log.Info("Skills retrieved successfully")
	return &proto.GetSkillsRes{
		Skill: skills,
	}, nil
}

func (s *UserService) AdDeleteSkill(ctx context.Context, req *proto.ADeleteSkillReq) (*proto.EmtpyRes, error) {
	err := s.reops.AdminDeleteSkill(uint(req.Id))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"skill_id": req.Id,
			"error":    err.Error(),
		}).Error("Failed to delete skill")
		return &proto.EmtpyRes{}, err
	}

	s.Log.WithField("skill_id", req.Id).Info("Skill deleted successfully")
	return &proto.EmtpyRes{}, nil
}

func (s *UserService) DeleteCategory(ctx context.Context, req *proto.DeleteCatReq) (*proto.EmtpyRes, error) {
	err := s.reops.AdminDeleteCategory(uint(req.Id))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"category_id": req.Id,
			"error":       err.Error(),
		}).Error("Failed to delete category")
		return &proto.EmtpyRes{}, err
	}

	s.Log.WithField("category_id", req.Id).Info("Category deleted successfully")
	return &proto.EmtpyRes{}, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *proto.EmtpyReq) (*proto.GetAllUserRes, error) {
	users, err := s.reops.GetAllUsers()
	if err != nil {
		s.Log.WithField("error", err.Error()).Error("Failed to retrieve users")
		return &proto.GetAllUserRes{}, err
	}

	s.Log.Info("All users retrieved successfully")
	return &proto.GetAllUserRes{
		Users: users,
	}, nil
}

func (s *UserService) AdminUserBlock(ctx context.Context, req *proto.BlockReq) (*proto.CommonRes, error) {
	user, err := s.reops.GetUserByID(uint(req.Id))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.Id,
			"error":   err.Error(),
		}).Error("Failed to find user")
		return &proto.CommonRes{}, err
	}
	if !user.IsActive {
		err := s.reops.UnBlockUser(user.ID)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id": req.Id,
				"action":  "active",
				"error":   err.Error(),
			}).Errorf("Failed to active user")
			return &proto.CommonRes{}, err
		}
	} else {
		err := s.reops.BlockUser(user.ID)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id": req.Id,
				"action":  "block",
				"error":   err.Error(),
			}).Errorf("Failed to block user")
			return &proto.CommonRes{}, err
		}
	}

	s.Log.WithFields(logrus.Fields{
		"user_id": req.Id,
	}).Info("User status updated successfully")
	return &proto.CommonRes{
		Message: "Updated Successfully",
		Status:  200,
	}, nil

}

func (s *UserService) GetUserEmail(ctx context.Context, req *proto.ProfileReq) (*proto.EmailRes, error) {
	fmt.Println("USerID", req.UserId)
	User, err := s.reops.GetUserByID(uint(req.UserId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"error":   err.Error(),
		}).Error("Failed to retrieve user")
		return nil, err
	}

	s.Log.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("User email retrieved successfully")
	return &proto.EmailRes{
		Email: User.Email,
	}, nil
}
