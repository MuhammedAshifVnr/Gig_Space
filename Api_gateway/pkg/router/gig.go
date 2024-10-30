package router

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func GigRouters(r fiber.Router, handller *handler.GigHandler) {
	r.Post("/add", middleware.Auth(Role), handller.CreateGig)
	r.Get("/user", middleware.Auth(Role), handller.GetGigByUserID)
	r.Put("/:id", middleware.Auth(Role), handller.UpdaeteGig)
	r.Delete("/:GigID", middleware.Auth(Role), handller.DeleteGig)
	r.Post("/order/:GigID", middleware.Auth(Role), handller.CreateOrder)
	r.Post("/quotes/:GigID", middleware.Auth(Role), handller.CreateQuote)
	r.Get("/quotes", middleware.Auth(Role), handller.GetAllQuote)
	r.Post("/custom", middleware.Auth(Role), handller.CreateCustomGig)
	r.Get("/offers", middleware.Auth(Role), handller.GetAllOffers)
	r.Post("/offers/:GigID", middleware.Auth(Role), handller.CreateOfferOrder)
	r.Get("/requests", middleware.Auth(Role), handller.GetAllOrdersRequest)
	r.Post("/:order_id/accept",middleware.Auth(Role),handller.AccepteOrder)
	r.Post("/:order_id/reject",middleware.Auth(Role),handller.RejectOrder)
}
