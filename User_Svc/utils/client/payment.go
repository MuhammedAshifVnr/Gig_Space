package client

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewPaymentClinet() proto.PaymentServiceClient {
	search, err := grpc.Dial(viper.GetString("PAYMENT_CONN"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}
	return proto.NewPaymentServiceClient(search)
}
