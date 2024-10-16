package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/broker"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/db"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/internal/service"
)

func InitializeService() *service.ChatService {
	DB := db.InitMongoDB()
	Repo := repo.NewChatRepository(DB)
	Rmq := broker.ConnectRabbitMQ()
	service:=service.NewChatService(Repo, Rmq)
	go service.ChatManger()
	return service
}
