package main

import (
	"log"
	"net"

	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/di"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config : %v", err)
	}
	service := di.InitializeService()
	server := grpc.NewServer()
	proto.RegisterChatServiceServer(server,service)
	lis,err:=net.Listen("tcp",viper.GetString("Port"))
	if err != nil {
		log.Fatalf("failed to run on the port %v : %v", viper.GetString("Port"), err)
	}

	log.Println("gRPC server is running on port ", viper.GetString("Port"))
	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
