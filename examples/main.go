package main

import (
	"context"
	"errors"
	"log"
	"time"

	o11yfiber "github.com/IsaacDSC/O11Y_fiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"go.opentelemetry.io/otel/attribute"

	oteltrace "go.opentelemetry.io/otel/trace"
)

const serviceName = "minimal-api"

func main() {
	tp := o11yfiber.StartTracing(o11yfiber.TracingConfig{
		EndpointCollector: "http://localhost:14268/api/traces",
		ServiceNameKey:    serviceName,
	})
	log.Fatal(o11yfiber.StartServerHttp(o11yfiber.SettingsHttp{
		ServiceNameMetrics: serviceName,
		TracerProvider:     tp,
		Handlers: []o11yfiber.Handler{
			{HandlerFunc: handleUser, Path: "/users/:id", Method: o11yfiber.GET},
			{HandlerFunc: handleError, Path: "/error", Method: o11yfiber.GET},
			{HandlerFunc: handleTimeout, Path: "/timeout", Method: o11yfiber.GET},
		},
		Middleware: []func(c *fiber.Ctx) error{
			cors.New(cors.Config{
				AllowOrigins: "https://gofiber.io, https://gofiber.net, *",
				AllowHeaders: "Origin, Content-Type, Accept",
				AllowMethods: "*",
			}),
			o11yfiber.MiddlewareIO("/error", "/api/live/ws"),
		},
		ServerPort: 3000,
	}))

}

func handleError(ctx *fiber.Ctx) error {
	return errors.New("abc")
}

func handleTimeout(c *fiber.Ctx) error {
	time.Sleep(time.Second * 20)
	return c.SendString("end timer")
}

func handleUser(c *fiber.Ctx) error {
	id := c.Params("id")
	name := getUser(c.UserContext(), id)
	return c.JSON(fiber.Map{"id": id, name: name})
}

func getUser(ctx context.Context, id string) string {
	_, span := o11yfiber.Span().Start(ctx, "getUser", oteltrace.WithAttributes(attribute.String("id", id)))
	defer span.End()
	if id == "123" {
		return "otelfiber tester"
	}
	return "unknown"
}
