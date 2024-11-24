package client

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewChatClient() proto.ChatServiceClient {
	Chatsvc, err := grpc.Dial(viper.GetString("CHAT_CONN"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Chat service: %v", err)
	}
	return proto.NewChatServiceClient(Chatsvc)
}
