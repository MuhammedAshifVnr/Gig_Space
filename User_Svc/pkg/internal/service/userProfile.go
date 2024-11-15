package service

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/upload"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/sirupsen/logrus"
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
			s.Log.WithFields(logrus.Fields{
				"user_id": req.UserId,
				"error":   err.Error(),
			}).Error("Failed to create bio")
			return &proto.CommonRes{}, err
		}
	} else {
		err := s.reops.UpdateBio(model.UserProfile{
			UserID: uint(req.UserId),
			Bio:    req.Description,
			Title:  req.Title,
		})
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id": req.UserId,
				"error":   err.Error(),
			}).Error("Failed to update bio")
			return &proto.CommonRes{}, err
		}
	}
	s.Log.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("Profile updated successfully")

	return &proto.CommonRes{
		Message: "Profile Updated.",
		Status:  200,
	}, nil
}

func (s *UserService) FreelacerAddSkill(ctx context.Context, req *proto.FreeAddSkillsReq) (*proto.CommonRes, error) {

	val, err := s.reops.GetSkill(req.SkillName)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"skill_name": req.SkillName,
			"error":      err.Error(),
		}).Error("Failed to find skill")
		return &proto.CommonRes{}, err
	}
	if val.ID == 0 {
		s.Log.WithFields(logrus.Fields{
			"skill_name": req.SkillName,
		}).Warn("Invalid skill provided")
		return &proto.CommonRes{
			Message: "Please enter valid skill.",
			Status:  400,
		}, nil
	}
	err = s.reops.FreelacerAddSkill(uint(req.UserId), val.ID, int(req.ProficiencyLevel))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id":           req.UserId,
			"skill_id":          val.ID,
			"proficiency_level": req.ProficiencyLevel,
			"error":             err.Error(),
		}).Error("Failed to save skill")
		return &proto.CommonRes{}, err
	}
	s.Log.WithFields(logrus.Fields{
		"user_id":  req.UserId,
		"skill_id": val.ID,
	}).Info("Skill updated successfully")

	return &proto.CommonRes{
		Message: "Skills Updated.",
		Status:  200,
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *proto.ProfileReq) (*proto.ProfileRes, error) {
	user, err := s.reops.GetProfile(uint(req.UserId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"error":   err.Error(),
		}).Error("Failed to find user profile") //
		return &proto.ProfileRes{}, err
	}
	skills, err := s.reops.GetSkillByuserID(uint(req.UserId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"error":   err.Error(),
		}).Error("Failed to find user skills")
		return &proto.ProfileRes{}, err
	}
	add := s.reops.GetAddress(uint(req.UserId))
	s.Log.WithFields(logrus.Fields{
		"user_id": req.UserId,
	}).Info("User profile retrieved successfully")

	return &proto.ProfileRes{
		Firstname:   user.FirstName,
		Lastname:    user.LastName,
		Email:       user.Email,
		Country:     user.Country,
		Phone:       user.Phone,
		Skill:       skills,
		Title:       user.Title,
		Description: user.Bio,
		Photo:       user.ProfilePhoto,
		Address: &proto.Address{
			City:     add.City,
			State:    add.State,
			District: add.District,
		},
	}, nil
}

func (s *UserService) DeleteSkill(ctx context.Context, req *proto.DeleteSkillRes) (*proto.CommonRes, error) {
	err := s.reops.DeleteSkillByID(uint(req.UserId), uint(req.SkillId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
		}).Info("User profile retrieved successfully")
		return &proto.CommonRes{}, err
	}
	s.Log.WithFields(logrus.Fields{
		"user_id":  req.UserId,
		"skill_id": req.SkillId,
	}).Info("Skill deleted successfully")

	return &proto.CommonRes{
		Message: "Skill Deleted Successfully",
		Status:  200,
	}, nil
}

func (s *UserService) UploadProfilePhoto(ctx context.Context, req *proto.UpProilePicReq) (*proto.CommonRes, error) {
	url, err := upload.UploadPhoto(s.s3, req.Pic, uint(req.UserId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"error":   err.Error(),
		}).Error("Failed to upload photo to S3")
		return &proto.CommonRes{}, err
	}
	profile, err := s.reops.GetPhoto(uint(req.UserId))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.UserId,
			"error":   err.Error(),
		}).Error("Failed to find existing photo")
		return &proto.CommonRes{}, err
	}
	if profile.ID == 0 {
		err := s.reops.CreatPhoto(url, uint(req.UserId))
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id": req.UserId,
				"url":     url,
				"error":   err.Error(),
			}).Error("Failed to save photo URL")
			return &proto.CommonRes{}, err
		}
	} else {
		err := s.reops.UpdatePhoto(url, uint(req.UserId))
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id": req.UserId,
				"url":     url,
				"error":   err.Error(),
			}).Error("Failed to update photo URL")
			return &proto.CommonRes{}, err
		}
	}
	s.Log.WithFields(logrus.Fields{
		"user_id": req.UserId,
		"url":     url,
	}).Info("Profile photo updated successfully")

	return &proto.CommonRes{
		Message: " Successfully Updated.",
		Status:  200,
	}, nil
}

func (s *UserService) UpdateAddress(ctx context.Context, req *proto.AddressReq) (*proto.CommonRes, error) {
	add := s.reops.GetAddress(uint(req.Id))
	if add.ID == 0 {
		err := s.reops.CreateAddress(req.State, req.District, req.City, int(req.Id))
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id":  req.Id,
				"state":    req.State,
				"district": req.District,
				"city":     req.City,
				"error":    err.Error(),
			}).Error("Failed to create address")
			return nil, err
		}
		s.Log.WithFields(logrus.Fields{
			"user_id":  req.Id,
			"state":    req.State,
			"district": req.District,
			"city":     req.City,
		}).Info("Address created successfully")
	} else {
		if req.City != "" {
			add.City = req.City
		}
		if req.District != "" {
			add.District = req.District
		}
		if req.State != "" {
			add.State = req.State
		}
		err := s.reops.UpdateAddress(add)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id":  req.Id,
				"state":    req.State,
				"district": req.District,
				"city":     req.City,
				"error":    err.Error(),
			}).Error("Failed to update address")
			return nil, err
		}
		s.Log.WithFields(logrus.Fields{
			"user_id":  req.Id,
			"state":    add.State,
			"district": add.District,
			"city":     add.City,
		}).Info("Address updated successfully")
	}
	return &proto.CommonRes{
		Message: "updated successfully",
		Status:  200,
	}, nil
}

func (s *UserService) UpdatRole(ctx context.Context, req *proto.RoleReq) (*proto.CommonRes, error) {
	user, err := s.reops.GetUserByID(uint(req.Id))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_id": req.Id,
			"error":   err.Error(),
		}).Error("Failed to find user")
		return nil, err
	}
	var role string
	if user.Role == "client" {
		role = "freelancer"
		if err := s.reops.RoleChange("freelancer", uint(req.Id)); err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id":  req.Id,
				"new_role": role,
				"error":    err.Error(),
			}).Error("Failed to change user role")
			return nil, err
		}
	} else {
		role = "clinet"
		if err := s.reops.RoleChange("client", uint(req.Id)); err != nil {
			s.Log.WithFields(logrus.Fields{
				"user_id":  req.Id,
				"new_role": role,
				"error":    err.Error(),
			}).Error("Failed to change user role")
			return nil, err
		}
	}
	s.Log.WithFields(logrus.Fields{
		"user_id": req.Id,
	}).Info("User role changed successfully")

	return &proto.CommonRes{
		Message: "Role Changed into " + role,
		Status:  200,
	}, nil
}
