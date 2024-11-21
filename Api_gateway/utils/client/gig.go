package client

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewGigClient() proto.GigServiceClient {
	gigsvc, err := grpc.Dial(viper.GetString("GIG_CONN"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Gig service: %v", err)
	}
	return proto.NewGigServiceClient(gigsvc)
}
