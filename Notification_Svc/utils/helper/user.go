package helper

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space/Notificaton_svc/utils/client"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

type NotificationService struct {
	userClient proto.UserServiceClient
	gigClient  proto.GigServiceClient
}

func NewNotificationService(userConn proto.UserServiceClient, gigConn proto.GigServiceClient) *NotificationService {
	return &NotificationService{
		userClient: userConn,
		gigClient:  gigConn,
	}
}

func GetUserEmail(userID uint) (string, error) {

	res, err := client.NewUserClinet().GetUserEmail(context.Background(), &proto.ProfileReq{
		UserId: uint32(userID),
	})

	return res.Email, err
}

func GetUserProfile(userID uint) (string, error) {
	res, err := client.NewUserClinet().GetUserProfile(context.Background(), &proto.ProfileReq{
		UserId: uint32(userID),
	})
	return res.Firstname + " " + res.Lastname, err
}

func GetOrderDetails(orderID string)(*proto.OrderDetail,error){
	res,err:=client.NewGigClinet().GetOrderByID(context.Background(),&proto.OrderByIDReq{
		OrderId: orderID,
	})
	return res,err
}