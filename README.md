# O11Y_fiber
Minimal API with fiber


### Start minimal API
*Use example*

```go
package main

import (
	"context"
	"errors"
	"log"

	o11yfiber "github.com/IsaacDSC/O11Y_fiber"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func main() {
	tp := o11yfiber.StartTracing(o11yfiber.TracingConfig{
		EndpointCollector: "http://localhost:14268/api/traces",
		ServiceNameKey:    "minimal-api",
	})

	log.Fatal(o11yfiber.StartServerHttp(o11yfiber.SettingsHttp{
		TracerProvider: tp,
		Handlers: []o11yfiber.Handler{
			{HandlerFunc: handleUser, Path: "/users/:id", Method: o11yfiber.GET},
			{HandlerFunc: handleError, Path: "/error", Method: o11yfiber.GET},
		},
		Middleware: []func(c *fiber.Ctx) error{
			o11yfiber.MiddlewareIO,
		},
		ServerPort: 3000,
	}))

}

func handleError(ctx *fiber.Ctx) error {
	return errors.New("abc")
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


```
