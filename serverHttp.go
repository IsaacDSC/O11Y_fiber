package o11yfiber

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Methods = string

const (
	POST   Methods = "POST"
	GET    Methods = "GET"
	PATCH  Methods = "PATCH"
	PUT    Methods = "PUT"
	DELETE Methods = "DELETE"
)

type HandlerFn = func(*fiber.Ctx) error
type responseWrapper = map[Methods]func(pattern Methods, handlerFn HandlerFn)

func wrapper(r *fiber.App) responseWrapper {
	routesWrap := make(responseWrapper)
	routesWrap[POST] = func(pattern string, handlerFn HandlerFn) {
		r.Post(pattern, handlerFn)
	}
	routesWrap[GET] = func(pattern string, handlerFn HandlerFn) {
		r.Get(pattern, handlerFn)
	}
	routesWrap[PATCH] = func(pattern string, handlerFn HandlerFn) {
		r.Patch(pattern, handlerFn)
	}
	routesWrap[PUT] = func(pattern string, handlerFn HandlerFn) {
		r.Put(pattern, handlerFn)
	}
	routesWrap[DELETE] = func(pattern string, handlerFn HandlerFn) {
		r.Delete(pattern, handlerFn)
	}
	return routesWrap
}

type Handler struct {
	HandlerFunc HandlerFn
	Path        string
	Method      string
}

type SettingsHttp struct {
	TracerProvider     *sdktrace.TracerProvider
	Handlers           []Handler
	Middleware         []func(c *fiber.Ctx) error
	ServerPort         int
	ServiceNameMetrics string
}

func StartServerHttp(config SettingsHttp) error {
	defer func() {
		if config.TracerProvider != nil {
			if err := config.TracerProvider.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}
	}()

	if config.ServerPort == 0 {
		config.ServerPort = 3000
	}

	if config.ServiceNameMetrics == "" {
		return errors.New(NOT_FOUND_SERVICE_NAME_METRICS)
	}

	server := fiber.New()

	prometheus := fiberprometheus.New(config.ServiceNameMetrics)
	prometheus.RegisterAt(server, "/metrics")

	server.Use(recover.New())
	server.Use(prometheus.Middleware)

	if config.TracerProvider != nil {
		server.Use(otelfiber.Middleware())
	}

	for i := range config.Middleware {
		server.Use(config.Middleware[i])
	}

	wrapper := wrapper(server)
	for i := range config.Handlers {
		method := config.Handlers[i].Method
		wrapper[method](config.Handlers[i].Path, config.Handlers[i].HandlerFunc)
	}

	return server.Listen(fmt.Sprintf(":%d", config.ServerPort))

}
