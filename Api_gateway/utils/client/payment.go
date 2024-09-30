package client

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewPaymentClient() proto.PaymentServiceClient {
	paymentSvc, err := grpc.Dial(viper.GetString("PayConn"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Payment service: %v", err)
	}
	return proto.NewPaymentServiceClient(paymentSvc)
}