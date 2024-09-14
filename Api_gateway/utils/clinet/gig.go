package clinet

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewGigClinet() proto.GigServiceClient{
	gigsvc, err := grpc.Dial(viper.GetString("GigConn"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Gig service: %v", err)
	}
	return proto.NewGigServiceClient(gigsvc)
}
