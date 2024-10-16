package broker

import (
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/spf13/viper"
)

func ConnectRabbitMQ() *amqp091.Connection {
	conn, err := amqp091.Dial(viper.GetString("AmqpUrl"))
	if err != nil {
		fmt.Println("error", err)
	}
	return conn
}
