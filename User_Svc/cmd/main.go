package main

import (
	"log"
	"net"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/di"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config : %v", err)
	}

	service := di.InitializeService(cfg)
	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server,service)
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to run on the port %v : %v", cfg.GRPCPort, err)
	}

	log.Println("gRPC server is running on port ", viper.GetString("Port"))
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
