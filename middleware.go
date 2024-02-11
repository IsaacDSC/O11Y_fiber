package o11yfiber

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func MiddlewareIO(ignorePaths ...string) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		l := NewLogger()

		var (
			correlationID = uuid.New().String()
			requestID     = uuid.New().String()
			path          = string(c.Request().URI().Path())
			endpoint      = c.Request().URI().String()
		)

		for i := range ignorePaths {
			if ignorePaths[i] == path {
				err := c.Next()
				if err != nil {
					return err
				}
				return nil
			}
		}

		ctx := context.WithValue(c.UserContext(), "correlation_id", correlationID)
		ctx = context.WithValue(ctx, "request_id", requestID)

		l.Log(ctx, slog.LevelInfo, "Request",
			"endpoint", endpoint,
			"path", path,
			"correlation_id", correlationID,
			"request_id", requestID,
		)

		err := c.Next()
		if err != nil {
			l.Log(ctx, slog.LevelError, "Response", "correlation_id", correlationID, "request_id", requestID)
			return err
		}

		l.Log(ctx, slog.LevelInfo, "Response",
			"correlation_id", correlationID,
			"request_id", requestID,
			"status_code", c.Response().StatusCode(),
			"body", string(c.Response().Body()),
			"endpoint", endpoint,
			"path", path,
		)

		return nil
	}
}
