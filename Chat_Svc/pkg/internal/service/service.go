package service

import (
	"context"
	"fmt"
	"log"

	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ChatService struct {
	Repo     *repo.ChatRepo
	AmqpConn *amqp.Connection
	proto.UnimplementedChatServiceServer
}

func NewChatService(repo *repo.ChatRepo, rmqConn *amqp.Connection) *ChatService {
	return &ChatService{
		Repo:     repo,
		AmqpConn: rmqConn,
	}
}

func (s *ChatService) ChatManger() {
	fmt.Println("-----")
	channel, err := s.AmqpConn.Channel()
	if err != nil {
		log.Println("Error to create channel: ", err)
	}
	defer channel.Close()

	qu, err := channel.QueueDeclare(
		"messageQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error to Dclare Queue: ", err)
	}
	message, err := channel.Consume(
		qu.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Error to Consume: ", err)
	}
	done := make(chan bool)
	go func() {
		for r := range message {
			log.Printf("Received a message: %s", r.Body)
			err := s.Repo.SaveMessages(r.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	<-done
}

func (s *ChatService) GetChat(ctx context.Context, req *proto.GetChatReq) (*proto.GetChatRes, error) {
	fmt.Println("res",req.RecipientId,req.SenderId)
	msg, err := s.Repo.GetMessages(req.SenderId, req.RecipientId)
	if err != nil {
		return nil, err
	}
	return &proto.GetChatRes{
		Chat: msg,
	}, nil
}
