package handler

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"time"

// 	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/helper"
// 	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/websocket/v2"
// 	amqp "github.com/rabbitmq/amqp091-go"
// )

// type ChatHandler struct {
// 	ChatClient proto.ChatServiceClient
// }

// func NewChatHandler(ChatConn proto.ChatServiceClient) *ChatHandler {
// 	return &ChatHandler{
// 		ChatClient: ChatConn,
// 	}
// }

// var (
// 	Controller = make(map[uint]*websocket.Conn)
// )

// func (h *ChatHandler) ChatConnection(c *websocket.Conn) {
// 	userID, ok := c.Locals("userID").(uint)
// 	fmt.Println("userID===", userID)
// 	if !ok {
// 		log.Println("userID is missing or not a uint")
// 		c.WriteJSON(fiber.Map{
// 			"error": "Unauthorized",
// 		})
// 		c.Close()
// 		return
// 	}
// 	Controller[userID] = c
// 	defer delete(Controller, userID)
// 	defer c.Close()

// 	for {
// 		_, msg, err := c.ReadMessage()
// 		if err != nil {
// 			c.WriteJSON(fiber.Map{
// 				"error": err.Error(),
// 			})
// 		}
// 		h.SendMessage(Controller, msg, userID)

// 	}
// }

// func (h *ChatHandler) SendMessage(user map[uint]*websocket.Conn, msg []byte, userID uint) {
// 	senderConn, ok := user[userID]
// 	var message helper.Message
// 	if err := json.Unmarshal([]byte(msg), &message); err != nil {
// 		if ok {
// 			senderConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
// 		}
// 		return
// 	}

// 	recipientConn, ok := user[message.RecipientID]
// 	if !ok {
// 		delete(user, message.RecipientID)
// 		err := h.RabbitmqSender(message)
// 		if err != nil {
// 			senderConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
// 		}
// 		return
// 	}
// 	err := h.RabbitmqSender(message)
// 	if err != nil {
// 		senderConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
// 	}

// 	err = recipientConn.WriteMessage(websocket.TextMessage, msg)
// 	if err != nil {
// 		senderConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
// 		delete(user, message.RecipientID)
// 	}
// }

// func (h *ChatHandler) RabbitmqSender(msg helper.Message) error {
// 	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
// 	if err != nil {
// 		return err
// 	}
// 	defer conn.Close()
// 	ch, err := conn.Channel()
// 	if err != nil {
// 		log.Fatalln("Failed to open a channel: ", err)
// 	}
// 	defer ch.Close()

// 	q, err := ch.QueueDeclare(
// 		"message",
// 		false,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return err
// 	}
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()
// 	msgbyte, err := json.Marshal(msg)
// 	if err != nil {
// 		return err
// 	}
// 	err = ch.PublishWithContext(
// 		ctx,
// 		"",
// 		q.Name,
// 		false,
// 		false,
// 		amqp.Publishing{
// 			ContentType: "text/plain",
// 			Body:        msgbyte,
// 		})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
