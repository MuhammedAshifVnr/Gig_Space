package client

import (
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewPaymentClinet() proto.PaymentServiceClient {
	Payment, err := grpc.Dial(viper.GetString("PAYMENT_CONN"), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to Payment service: %v", err)
	}
	return proto.NewPaymentServiceClient(Payment)
}
