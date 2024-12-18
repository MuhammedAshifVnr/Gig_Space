package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/logger"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/jwt"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/otp"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	reops         repo.RepoInter
	PaymentClient proto.PaymentServiceClient
	s3            *s3.S3
	Log           *logrus.Logger
	proto.UnimplementedUserServiceServer
}

func NewUserService(repo repo.RepoInter, S3 *s3.S3, payClient proto.PaymentServiceClient) *UserService {
	return &UserService{
		PaymentClient: payClient,
		reops:         repo,
		s3:            S3,
		Log:           logger.Log,
	}
}

func (s *UserService) UserSignup(ctx context.Context, req *proto.SignupReq) (*proto.SignupRes, error) {
	// if req.Email == "" {
	// 	return &proto.SignupRes{
	// 		Message: "email filed can't be empty",
	// 		Status:  http.StatusBadRequest,
	// 	}, nil
	// }
	// if req.Password == "" {
	// 	return &proto.SignupRes{
	// 		Message: "password filed can't be empty",
	// 		Status:  http.StatusBadRequest,
	// 	}, nil
	// }
	// if req.Firstname == "" {
	// 	return &proto.SignupRes{
	// 		Message: "first name filed can't be empty",
	// 		Status:  http.StatusBadRequest,
	// 	}, nil
	// }
	// if req.Lastname == "" {
	// 	return &proto.SignupRes{
	// 		Message: "last name filed can't be empty",
	// 		Status:  http.StatusBadRequest,
	// 	}, nil
	// }
	// if req.Phone == "" {
	// 	return &proto.SignupRes{
	// 		Message: "please enter a valied phone number",
	// 		Status:  http.StatusBadRequest,
	// 	}, nil
	// }

	err := s.reops.CheckingExist(req.Email, req.Phone)
	if err != nil {
		s.Log.WithError(err).Error("Failed to Checking Existing Status")
		return &proto.SignupRes{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}, nil
	}

	Otp, err := otp.SendOtp(req.Email, string(req.Firstname+" "+req.Lastname))
	if err != nil {
		s.Log.WithError(err).Errorf("Failed to Send Otp: %v", err)
		return &proto.SignupRes{
			Error:  err.Error(),
			Status: http.StatusForbidden,
		}, nil
	}

	err = s.reops.SignupData(model.User{
		FirstName: req.Firstname,
		LastName:  req.Lastname,
		Email:     req.Email,
		Password:  req.Password,
		Role:      req.Role,
		Phone:     req.Phone,
		Country:   req.Country,
		IsActive:  true,
	}, Otp)
	if err != nil {
		s.Log.WithError(err).Errorf("Failed to save user data: %v", err)
		return &proto.SignupRes{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		}, nil
	}

	s.Log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("User signed up successfully")

	return &proto.SignupRes{
		Message: "Verifcaton link sent to email.Verify to get access",
		Status:  200,
	}, nil

}

func (s *UserService) VerifyingEmail(ctx context.Context, req *proto.VerifyReq) (*proto.VerifyRes, error) {
	val, err := s.reops.VerifyingEmail(req.Otp, req.Email)
	if err == redis.Nil {
		s.Log.Warn("Verification link expired")
		return &proto.VerifyRes{}, errors.New("this link was expired")
	}
	var user model.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		s.Log.WithError(err).Error("Failed to unmarshal user data")
		return &proto.VerifyRes{}, errors.New("Could not unmarshal user: " + err.Error())
	}
	err = s.reops.CreateUser(user)
	if err != nil {
		s.Log.WithError(err).Error("Failed to create user") // Log error with details
		return &proto.VerifyRes{}, err
	}
	s.reops.DeleteOtp(req.Otp)
	s.Log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("User verified successfully")

	return &proto.VerifyRes{
		Message: "User Verifed Successfully.Go to the Login.",
		Status:  200,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginRes, error) {
	fmt.Println(req.Email)
	user, err := s.reops.GetUser(req.Email)
	if err != nil {
		s.Log.WithError(err).Errorf("Failed to find user: %v", req.Email)
		return &proto.LoginRes{
			Message: "User not fount",
			Status:  http.StatusNotFound,
			Error:   err.Error(),
		}, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"email": req.Email,
		}).Warn("Invalid username or password")
		return &proto.LoginRes{
			Status:  http.StatusUnauthorized,
			Message: "Invalid Username or Password.",
			Error:   err.Error(),
		}, nil
	}
	if !user.IsActive {
		s.Log.WithFields(logrus.Fields{
			"email": req.Email,
		}).Warn("User is blocked")
		return &proto.LoginRes{
			Status:  http.StatusUnauthorized,
			Message: "User is Bolcked",
		}, nil
	}
	res, err := s.PaymentClient.GetSubDetails(context.Background(), &proto.GetSubReq{UserId: uint32(user.ID)})
	if err != nil {
		return nil, err
	}
	token, err := jwt.GenerateJwtToken(user.Email, user.ID, "user", res.ExpiryTime)
	if err != nil {
		s.Log.WithError(err).Error("Error during JWT token generation")
		return &proto.LoginRes{
			Message: "Error form jwt creation ",
			Status:  404,
			Error:   err.Error(),
		}, nil
	}

	s.Log.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info("User Login successfully")

	return &proto.LoginRes{
		Message: "login successful",
		Status:  200,
		Token:   token,
	}, nil
}

func (s *UserService) ForgotPassword(ctx context.Context, req *proto.FP_Req) (*proto.CommonRes, error) {
	user, err := s.reops.GetUser(req.Email)
	if err != nil {
		s.Log.WithError(err).Errorf("Failed to find user: %v", req.Email)
		return &proto.CommonRes{}, fmt.Errorf("user not fount")
	}
	otp, err := otp.ForgotOtp(user.Email, user.FirstName)
	if err != nil {
		s.Log.WithError(err).Error("Failed to generate forgot password OTP")
		return &proto.CommonRes{}, err
	}
	err = s.reops.SignupData(model.User{Email: user.Email}, otp)
	if err != nil {
		s.Log.WithError(err).Error("Failed to save forgot password data")
		return &proto.CommonRes{}, err
	}

	s.Log.WithFields(logrus.Fields{
		"user_id": req.Email,
	}).Info("Reset OTP Successfully Sended")

	return &proto.CommonRes{
		Message: "Reset OTP Successfully Sended",
		Status:  200,
	}, nil
}

func (s *UserService) ResetPassword(ctx context.Context, req *proto.ResetPwdReq) (*proto.CommonRes, error) {
	val, err := s.reops.VerifyingEmail(req.Otp, "")
	if err == redis.Nil {
		s.Log.WithFields(logrus.Fields{
			"otp": req.Otp,
		}).Error("Invalid or expired OTP")
		return &proto.CommonRes{}, errors.New("this OTP was invalid or expired")
	}
	var user model.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"otp":   req.Otp,
			"error": err.Error(),
		}).Error("Failed to unmarshal user")
		return &proto.CommonRes{}, errors.New("Could not unmarshal user: " + err.Error())
	}
	err = s.reops.ResetPassword(user.Email, req.Password)
	if err != nil {
		s.Log.WithFields(logrus.Fields{
			"user_email": user.Email,
			"error":      err.Error(),
		}).Error("Failed to reset password")
		return &proto.CommonRes{}, err
	}
	s.Log.WithFields(logrus.Fields{
		"user_email": user.Email,
	}).Info("Password updated successfully")

	return &proto.CommonRes{
		Message: "Password Updated",
		Status:  200,
	}, nil
}
