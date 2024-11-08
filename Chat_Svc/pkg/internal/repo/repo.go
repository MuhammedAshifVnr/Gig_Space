package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/Chat_Svc/pkg/model"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatRepo struct {
	DB  *mongo.Collection
	RDB *redis.Client
}

func NewChatRepository(conn *mongo.Collection, RedisConn *redis.Client) *ChatRepo {
	return &ChatRepo{DB: conn, RDB: RedisConn}
}

func (r *ChatRepo) SaveMessages(msg []byte) error {
	var message model.Message
	err := json.Unmarshal(msg, &message)
	if err != nil {
		return err
	}
	message.CreatedAt = time.Now()
	_, err = r.DB.InsertOne(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}

func (r *ChatRepo) GetMessages(senderID, recipientID uint32) ([]*proto.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"senderid": bson.M{"$in": bson.A{senderID, recipientID}}, "recipientid": bson.M{"$in": bson.A{senderID, recipientID}}}

	opts := options.Find().SetSort(bson.D{{"created_at", 1}})

	cursor, err := r.DB.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*proto.Message
	for cursor.Next(ctx) {
		var msg model.Message
		if err := cursor.Decode(&msg); err != nil {
			return nil, err
		}

		protoMsg := &proto.Message{
			SenderId:    strconv.Itoa(int(msg.SenderID)),
			RecipientId: strconv.Itoa(int(msg.RecipientID)),
			TextMessage: msg.MessageText,
			Time:        msg.CreatedAt.String(),
		}
		messages = append(messages, protoMsg)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *ChatRepo) StoreNotification(senderID, recipientID int32) bool {
	ctx := context.Background()
	key := fmt.Sprintf("notify:%d:%d", senderID, recipientID)

	exists, err := r.RDB.Exists(ctx, key).Result()
	if err != nil {
		fmt.Println("Error accessing Redis:", err)
		return false
	}

	if exists == 0 {
		r.RDB.Set(ctx, key, "1", time.Hour)
		return true
	}

	return false
}
