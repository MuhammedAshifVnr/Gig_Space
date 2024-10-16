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
}
