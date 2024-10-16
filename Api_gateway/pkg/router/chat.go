package router

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func ChatRouter(r fiber.Router, handler *handler.MessagingHandler) {
	r.Get("", middleware.AuthChat(), websocket.New(handler.OpenChat))
	r.Get("/messages/:receiverID", middleware.Auth(Role), handler.GetChat)
}
