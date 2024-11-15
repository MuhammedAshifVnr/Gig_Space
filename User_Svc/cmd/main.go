package main

import (
	"net"
	"os"

	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/di"
	"github.com/MuhammedAshifVnr/Gig_Space/User_svc/pkg/logger"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"google.golang.org/grpc"
)

func main() {

	logger.Init()
	log := logger.Log
	log.Info("Starting application...")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to load configuration")
		os.Exit(1)
	}
	log.Info("Configuration loaded successfully")

	service := di.InitializeService(cfg)
	log.Info("Service dependencies initialized")
	server := grpc.NewServer()
	proto.RegisterUserServiceServer(server, service)
	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.WithFields(logrus.Fields{
			"port":  viper.GetString("Port"),
			"error": err,
		}).Fatal("Failed to listen on port")
		os.Exit(1)
	}

	log.WithFields(logrus.Fields{
		"port": viper.GetString("Port"),
	}).Info("gRPC server is running")
	if err := server.Serve(lis); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to serve gRPC server")
		os.Exit(1)
	}
}
