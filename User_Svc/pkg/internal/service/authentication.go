package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/jwt"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/utils/otp"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	reops repo.RepoInter
	proto.UnimplementedUserServiceServer
}

func NewUserService(repo repo.RepoInter) *UserService {
	return &UserService{
		reops: repo,
	}
}

func (s *UserService) UserSignup(ctx context.Context, req *proto.SignupReq) (*proto.SignupRes, error) {
	if req.Email == "" {
		return &proto.SignupRes{
			Message: "email filed can't be empty",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if req.Password == "" {
		return &proto.SignupRes{
			Message: "password filed can't be empty",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if req.Firstname == "" {
		return &proto.SignupRes{
			Message: "first name filed can't be empty",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if req.Lastname == "" {
		return &proto.SignupRes{
			Message: "last name filed can't be empty",
			Status:  http.StatusBadRequest,
		}, nil
	}
	if req.Phone == "" {
		return &proto.SignupRes{
			Message: "please enter a valied phone number",
			Status:  http.StatusBadRequest,
		}, nil
	}

	// err := s.reops.CheckingExist(req.Email, req.Phone)
	// if err != nil {
	// 	return &proto.SignupRes{
	// 		Message: err.Error(),
	// 		Status:  http.StatusBadRequest,
	// 	}, nil
	// }

	Otp, err := otp.SendOtp(req.Email, string(req.Firstname+" "+req.Lastname))
	if err != nil {
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
	}, Otp)
	if err != nil {
		return &proto.SignupRes{
			Error:  err.Error(),
			Status: http.StatusBadRequest,
		}, nil
	}
	return &proto.SignupRes{
		Message: "Verifcaton link sent to email.Verify to get access",
		Status:  200,
	}, nil

}

func (s *UserService) VerifyingEmail(ctx context.Context, req *proto.VerifyReq) (*proto.VerifyRes, error) {
	val, err := s.reops.VerifyingEmail(req.Otp, req.Email)
	if err == redis.Nil {
		return &proto.VerifyRes{}, errors.New("this link was expired")
	}
	var user model.User
	err = json.Unmarshal([]byte(val), &user)
	if err != nil {
		return &proto.VerifyRes{}, errors.New("Could not unmarshal user: " + err.Error())
	}
	err = s.reops.CreateUser(user)
	if err != nil {
		return &proto.VerifyRes{}, err
	}
	return &proto.VerifyRes{
		Message: "User Verifed Successfully.Go to the Login.",
		Status:  200,
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *proto.LoginReq) (*proto.LoginRes, error) {
	fmt.Println(req.Email)
	user, err := s.reops.GetUser(req.Email)
	if err != nil {
		return &proto.LoginRes{
			Message: "User not fount",
			Status:  http.StatusNotFound,
			Error:   err.Error(),
		}, nil
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &proto.LoginRes{
			Status:  http.StatusUnauthorized,
			Message: "Invalid Username or Password.",
			Error:   err.Error(),
		}, nil
	}
	token, err := jwt.GenerateJwtToken(user.Email,user.ID,"user")
	if err != nil {
		return &proto.LoginRes{
			Message: "Error form jwt creation ",
			Status:  404,
			Error:   err.Error(),
		}, nil
	}
	return &proto.LoginRes{
		Message: "login successful",
		Status:  200,
		Token:   token,
	}, nil
}

func (s *UserService)ForgetPassword(ctx context.Context,)(){

}
