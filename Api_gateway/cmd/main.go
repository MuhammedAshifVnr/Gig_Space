package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/MuhammedAshifVnr/Gig_Space/api_gateway/docs"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/config"
	"github.com/MuhammedAshifVnr/Gig_Space/api_gateway/pkg/di"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

// @title API Gateway Swagger
// @version 1.0
// @description This is the API Gateway for the Flexi Worke project

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal("---", err)
	}

	app := fiber.New()

	app.Use(logger.New())

	app.Use(limiter.New(limiter.Config{
		Max:        59,
		Expiration: 60 * time.Second,
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		// Render the HTML file from the 'views' folder
		return c.Render("temp/payment.html", fiber.Map{
			"Title": "Fiber HTML Example",
		})
	})

	app.Get("/swagger/*", swagger.HandlerDefault)
	app.Static("/temp", "./temp")
	di.InitializeRoutes(app)

	http.Handle("/metrics", promhttp.Handler())

	// Start Prometheus server
	go func() {
		log.Println("Starting Prometheus metrics server on port 9090")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Fatal("Prometheus metrics server stopped:", err)
		}
	}()

	fmt.Println("Port", viper.GetString("PORT"))
	if err := app.Listen(viper.GetString("PORT")); err != nil {
		fmt.Println("stoped", err)
	}

}
