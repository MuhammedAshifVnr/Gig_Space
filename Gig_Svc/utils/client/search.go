package client

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewSearchClinet() proto.SearchServiceClient {
	search, err := grpc.Dial(viper.GetString("SEARCH_CONN"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Search service: %v", err)
	}
	return proto.NewSearchServiceClient(search)
}
