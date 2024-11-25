package monitoring

import (
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

func Register(app *fiber.App) {
	// Create a new Prometheus registry
	registry := prometheus.NewRegistry()

	// Create a Prometheus HTTP handler
	prometheusHandler := fasthttpadaptor.NewFastHTTPHandler(promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// Define HTTP requests counter metric
	requestCounter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests received",
		},
		[]string{"method", "path"},
	)
	registry.MustRegister(requestCounter)

	serviceHealth := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "service_health",
			Help: "Indicates the health of the service (1 = up, 0 = down)",
		},
		[]string{"service"},
	)
	registry.MustRegister(serviceHealth)

	serviceHealth.WithLabelValues("api_gateway").Set(1)

	app.Use(func(c *fiber.Ctx) error {
		requestCounter.WithLabelValues(c.Method(), c.Path()).Inc()
		return c.Next()
	})

	app.Get("/metrics", func(c *fiber.Ctx) error {
		prometheusHandler(c.Context())
		return nil
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	app.Get("/simulate_down", func(c *fiber.Ctx) error {
		serviceHealth.WithLabelValues("api_gateway").Set(0)
		return c.JSON(fiber.Map{"status": "Service marked as DOWN"})
	})

	app.Get("/simulate_up", func(c *fiber.Ctx) error {
		serviceHealth.WithLabelValues("api_gateway").Set(1)
		return c.JSON(fiber.Map{"status": "Service marked as UP"})
	})
}
