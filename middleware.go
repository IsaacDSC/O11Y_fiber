package o11yfiber

import (
	"context"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func MiddlewareIO(c *fiber.Ctx) error {
	l := NewLogger()
	var (
		correlationID = uuid.New().String()
		requestID     = uuid.New().String()
	)
	ctx := context.WithValue(c.UserContext(), "correlation_id", correlationID)
	ctx = context.WithValue(ctx, "request_id", requestID)

	l.Log(ctx, slog.LevelInfo, "Request",
		"endpoint", c.Request().URI().String(),
		"path", string(c.Request().URI().RequestURI()),
		"correlation_id", correlationID,
		"request_id", requestID,
	)

	err := c.Next()
	if err != nil {
		l.Log(ctx, slog.LevelError, "Response", "correlation_id", correlationID, "request_id", requestID)
		return err
	}

	l.Log(ctx, slog.LevelError, "Response",
		"correlation_id", correlationID,
		"request_id", requestID,
		"status_code", c.Response().StatusCode(),
		"body", string(c.Response().Body()),
		"endpoint", c.Request().URI().String(),
		"path", string(c.Request().URI().RequestURI()),
	)

	return nil
}
