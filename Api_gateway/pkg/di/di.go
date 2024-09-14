package di

import (
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/handler"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/router"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/utils/clinet"
	"github.com/gofiber/fiber/v2"
)

func InitializeRoutes(app *fiber.App) {

	user := handler.NewUserHandler(clinet.NewUserClinet())
	admin:=handler.NewUserHandler(clinet.NewUserClinet())
	gig:=handler.NewGigHandler(clinet.NewGigClinet())
	router.UserProfile(app.Group("/user"), user)
	router.AdminRouters(app.Group("/admin"), admin)
	router.GigRouters(app.Group("/gig"),gig)
}
