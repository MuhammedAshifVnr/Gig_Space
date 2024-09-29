package router

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/gofiber/fiber/v2"
)

func SearchRouter(r fiber.Router, handller *handler.SearchHandler) {
	r.Get("/gigs", handller.SearchGig)
}
