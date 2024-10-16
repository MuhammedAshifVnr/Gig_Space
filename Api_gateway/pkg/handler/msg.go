package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/helper"
	"github.com/MuhammedAshifVnr/Gig_Space_Proto/proto"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rabbitmq/amqp091-go"
)

type MessagingHandler struct {
	MsgClient proto.ChatServiceClient
}

func NewChatHandler(client proto.ChatServiceClient) *MessagingHandler {
	return &MessagingHandler{
		MsgClient: client,
	}
}

var (
	ActiveConnections = make(map[int32]*websocket.Conn)
)

func (h *MessagingHandler) OpenChat(c *websocket.Conn) {
	userID, exists := c.Locals("userID").(uint)
	if !exists {
		log.Println("Invalid or missing userID")
		c.WriteJSON(fiber.Map{
			"error": "Access Denied",
		})
		_ = c.Close()
		return
	}
	ActiveConnections[int32(userID)] = c
	defer func() {
		delete(ActiveConnections, int32(userID))
		c.Close()
	}()

	for {
		_, incomingMsg, err := c.ReadMessage()
		if err != nil {
			c.WriteJSON(fiber.Map{
				"error": err.Error(),
			})
			continue
		}
		h.DispatchMessage(ActiveConnections, incomingMsg, userID)
	}
}

func (h *MessagingHandler) DispatchMessage(users map[int32]*websocket.Conn, msg []byte, userID uint) {
	senderConn, exists := users[int32(userID)]
	var msgData helper.Message
	if err := json.Unmarshal(msg, &msgData); err != nil {
		if exists {
			senderConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}
		return
	}
	msgData.SenderID = int32(userID)
	recipientConn, isOnline := users[msgData.RecipientID]
	if !isOnline {
		delete(users, msgData.RecipientID)
		err := h.PublishToQueue(msgData)
		if err != nil {
			senderConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}
		return
	}

	if err := h.PublishToQueue(msgData); err != nil {
		senderConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
	}
	if err := recipientConn.WriteMessage(websocket.TextMessage, msg); err != nil {
		senderConn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		delete(users, msgData.RecipientID)
	}
}

func (h *MessagingHandler) PublishToQueue(msg helper.Message) error {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Println("Error in opening RabbitMQ channel:", err)
		return err
	}
	defer conn.Channel()

	queue, err := channel.QueueDeclare(
		"messageQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = channel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        msgBytes,
		})
	if err != nil {
		return err
	}
	return nil
}

// @Summary      Get Chat Messages
// @Description  Retrieve the chat history between the current authenticated user and a specified recipient.
// @Tags         Chat
// @Accept       json
// @Produce      json
// @Param        receiverID  path      int  true  "Recipient User ID"
// @Router       /chat/messages/{receiverID} [get]
func (h *MessagingHandler) GetChat(c *fiber.Ctx) error {
	userID, _ := c.Locals("userID").(uint)
	fmt.Println("user", userID)
	reciptent, _ := strconv.Atoi(c.Params("receiverID"))
	fmt.Println("user", reciptent)
	res, err := h.MsgClient.GetChat(context.Background(), &proto.GetChatReq{
		SenderId:    uint32(userID),
		RecipientId: uint32(reciptent),
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(200).JSON(res)
}
