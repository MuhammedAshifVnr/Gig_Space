package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/broker"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/db"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/internal/service"
	"github.com/spf13/viper"
)

func InitializeService() *service.ChatService {
	DB := db.InitMongoDB()
	RDB := db.InitRedisClient()
	Repo := repo.NewChatRepository(DB, RDB)
	Rmq := broker.ConnectRabbitMQ()
	Kafka := config.InitKafkaWriters([]string{viper.GetString("Broker")}, []string{viper.GetString("OfflineTopic")})
	service := service.NewChatService(Repo, Rmq,Kafka)
	go service.ChatManger()
	return service
}
