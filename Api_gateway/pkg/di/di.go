package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/router"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/client"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {

	user := handler.NewUserHandler(client.NewUserClient())
	admin := handler.NewUserHandler(client.NewUserClient())
	gig := handler.NewGigHandler(client.NewGigClient())
	search := handler.NewSearchHandler(client.NewSearchClient())
	payment:=handler.NewPaymentHnadller(client.NewPaymentClient())
	chat:=handler.NewChatHandler(client.NewChatClient())
	router.SearchRouter(app.Group("/search"), search)
	router.UserProfile(app.Group("/user"), user)
	router.AdminRouters(app.Group("/admin"), admin,payment)
	router.GigRouters(app.Group("/gig"), gig)
	router.PaymentRouters(app.Group("/payments"),payment)
	router.ChatRouter(app.Group("/chat"),chat)
}
