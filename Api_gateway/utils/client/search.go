package client

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewSearchClient() proto.SearchServiceClient {
	SearchClient, err := grpc.Dial(viper.GetString("SearchConn"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Search service: %v", err)
	}
	return proto.NewSearchServiceClient(SearchClient)
}
