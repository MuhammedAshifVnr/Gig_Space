package helper

import (
	"context"
	"fmt"

	"github.com/MuhammedAshifVnr/Gig_Space/Notificaton_svc/utils/client"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
)

type NotificationService struct {
	userClient proto.UserServiceClient
	gigClient proto.GigServiceClient
}

func NewNotificationService(userConn proto.UserServiceClient,gigConn proto.GigServiceClient) *NotificationService {
	return &NotificationService{
		userClient: userConn,
		gigClient: gigConn,
	}
}

func GetUserEmail(userID uint) (string, error) {
	fmt.Println("user", userID)
	res, err := client.NewUserClinet().GetUserEmail(context.Background(), &proto.ProfileReq{
		UserId: uint32(userID),
	})
	return res.Email, err
}
