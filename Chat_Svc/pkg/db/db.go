package db

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitMongoDB() *mongo.Collection {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(viper.GetString("Mongo_URI")))
	if err != nil {
		log.Fatalf("Failed to connect DB :%v", err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to connect DB :%v", err)
	}
	return client.Database("datas").Collection("messages")
}

func InitRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
