package router

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

func PaymentRouters(r fiber.Router, handller *handler.PaymentHandler) {
	r.Post("/subscription/:PlanID", middleware.PaymentAuth(), handller.CreateSubscriptionPayment)
	r.Post("/razorpay-payment", handller.UpdatePaymentStatus)
	r.Post("/subscription-renew/:PlanID", middleware.PaymentAuth(), handller.RenewSubscription)
	r.Get("/", handller.RenderPaymentPage)
}
