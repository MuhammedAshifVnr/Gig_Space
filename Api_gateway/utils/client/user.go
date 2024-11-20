package client

import (
	"fmt"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewUserClient() proto.UserServiceClient {
	usersvc, err := grpc.Dial(viper.GetString("USER_CONN"), grpc.WithInsecure())
	if err != nil {
		fmt.Println("User Port: ", viper.GetString("USER_CONN"))
		log.Fatalf("failed to connect to User service: %v", err)
	}
	return proto.NewUserServiceClient(usersvc)
}
