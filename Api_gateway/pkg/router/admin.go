package router

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

const role ="AdminToken"

func AdminRouters(r fiber.Router, handller *handler.UserHandler){
	r.Post("/login",handller.AdminLogin)
	r.Post("/category",middleware.Auth(role),handller.AddCategory)
	r.Post("/skill",middleware.Auth(role),handller.AddSkill)
	r.Get("/category",middleware.Auth(role),handller.GetCategory)
	r.Get("/skills",middleware.Auth(role),handller.GetSkills)
	r.Delete("/skills/:skillID",middleware.Auth(role),handller.AdDeleteSkill)
	r.Delete("/category/:CatID",middleware.Auth(role),handller.AdDeleteCategory)
	r.Get("/users",middleware.Auth(role),handller.GetAllUsers)
}
