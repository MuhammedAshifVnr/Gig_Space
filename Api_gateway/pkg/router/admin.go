package router

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

const role = "AdminToken"

func AdminRouters(r fiber.Router, handller *handler.UserHandler, payHandller *handler.PaymentHandler, gigHandler *handler.GigHandler) {
	r.Post("/login", handller.AdminLogin)
	r.Post("/category", middleware.Auth(role), handller.AddCategory)
	r.Post("/skill", middleware.Auth(role), handller.AddSkill)
	r.Get("/category", middleware.Auth(role), handller.GetCategory)
	r.Get("/skills", middleware.Auth(role), handller.GetSkills)
	r.Delete("/skills/:skillID", middleware.Auth(role), handller.AdDeleteSkill)
	r.Delete("/category/:CatID", middleware.Auth(role), handller.AdDeleteCategory)
	r.Get("/users", middleware.Auth(role), handller.GetAllUsers)
	r.Post("/block/:userID", middleware.Auth(role), handller.UserBlock)
	r.Post("/create-plan", middleware.Auth(role), payHandller.CreatePlan)
	r.Get("/plans", middleware.Auth(role), payHandller.GetPlans)
	r.Delete("/plan/:PlanID", middleware.Auth(role), payHandller.DeletPlan)
	r.Post("/logout", middleware.Auth(role), handller.AdminLogout)
	r.Get("/orders/refund", middleware.Auth(Role), gigHandler.GetAllRefundOrders)
	r.Post("/orders/refund/:OrderID", middleware.Auth(Role), gigHandler.Refund)
	r.Get("/orders/completed", middleware.Auth(Role), gigHandler.GetAllCompletedOrders)
	r.Post("/orders/complet/:OrderID", middleware.Auth(Role), gigHandler.Payment)

}
