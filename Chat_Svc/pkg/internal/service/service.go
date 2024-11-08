package service

import (
	"context"

	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/internal/repo"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ChatService struct {
	Repo        *repo.ChatRepo
	AmqpConn    *amqp.Connection
	kafkaWriter map[string]*kafka.Writer
	proto.UnimplementedChatServiceServer
}

func NewChatService(repo *repo.ChatRepo, rmqConn *amqp.Connection, kafka map[string]*kafka.Writer) *ChatService {
	return &ChatService{
		Repo:        repo,
		AmqpConn:    rmqConn,
		kafkaWriter: kafka,
	}
}

func (s *ChatService) ChatManger() {
	channel, err := s.AmqpConn.Channel()
	if err != nil {
		logrus.WithError(err).Error("Error creating channel")
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
		logrus.WithError(err).Error("Error declaring queue")
		//
		return
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
		logrus.WithError(err).Error("Error consuming messages from queue")
		//
		return
	}
	done := make(chan bool)
	go func() {
		for r := range message {
			logrus.WithField("message", string(r.Body)).Info("Received a message")
			err := s.Repo.SaveMessages(r.Body)
			if err != nil {
				logrus.WithError(err).Error("Failed to save message to repository")
			}
		}
	}()
	<-done
}

func (s *ChatService) GetChat(ctx context.Context, req *proto.GetChatReq) (*proto.GetChatRes, error) {

	msg, err := s.Repo.GetMessages(req.SenderId, req.RecipientId)
	if err != nil {
		logrus.WithError(err).Error("Error fetching messages from repository")
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"sender_id":    req.SenderId,
		"recipient_id": req.RecipientId,
	}).Info("Successfully fetched chat messages")

	return &proto.GetChatRes{
		Chat: msg,
	}, nil
}

func (s *ChatService) SendOflineNotification(ctx context.Context, req *proto.NotifiyReq) (*proto.ChatCommonRes, error) {
	sent := s.Repo.StoreNotification(int32(req.SenderId), int32(req.RecipientId))
	if sent {
		err := helper.SendNotification(ctx, model.ChatEvent{
			SenderID:    int32(req.SenderId),
			RecipientID: int32(req.RecipientId),
			Event:       "Offline",
		}, viper.GetString("OfflineTopic"), s.kafkaWriter[viper.GetString("OfflineTopic")])
		if err != nil {
			logrus.WithError(err).Error("Failed to send offline notification")
			return nil, err
		}
	}
	
	logrus.WithFields(logrus.Fields{
		"sender_id":    req.SenderId,
		"recipient_id": req.RecipientId,
	}).Info("Successfully stored offline notification")
	
	return nil, nil
}
