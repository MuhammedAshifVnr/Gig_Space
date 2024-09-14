package clinet

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewUserClinet() proto.UserServiceClient {
	usersvc, err := grpc.Dial(viper.GetString("UserConn"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to User service: %v", err)
	}
	return proto.NewUserServiceClient(usersvc)
}
