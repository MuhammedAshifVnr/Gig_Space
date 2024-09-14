package router

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/middleware"
	"github.com/gofiber/fiber/v2"
)

const Role ="UserToken"

func UserProfile(r fiber.Router,handller *handler.UserHandler){
	r.Post("/signup",handller.Signup)
	r.Get("/verify",handller.VerifyingEmail)
	r.Post("/login",handller.Login)
	r.Post("/skills",middleware.Auth(Role),handller.FreeAddSkills)
	r.Put("/profile",middleware.Auth(Role),handller.UpdateProfile)
	r.Get("/profile",middleware.Auth(Role),handller.GetProfile)
	r.Delete("/skill/:Skill",middleware.Auth(Role),handller.DeleteSkill)
	r.Post("/profile-photo",middleware.Auth(Role),handller.UploadProfilePhoto)
}


