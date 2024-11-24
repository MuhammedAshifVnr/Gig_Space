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
	r.Post("/wallet/create", middleware.Auth(Role), handller.CreateWallet)
	r.Post("/wallet", middleware.Auth(Role), handller.GetWallet)
	r.Post("/bank", middleware.Auth(Role), handller.AddBankAccount)
	r.Post("/withdrawal", middleware.Auth(Role), handller.Withdrawal)
	r.Post("/wallet/change-pin", middleware.Auth(Role), handller.ChangeWalletPin)
	r.Post("/wallet/forgot-pin", middleware.Auth(Role), handller.ForgotWalletPin)
	r.Post("/wallet/reset-pin", middleware.Auth(Role), handller.ResetWalletPin)
	r.Post("/webhook", handller.UpdateWebhook)
}
