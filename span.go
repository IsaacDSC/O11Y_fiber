package o11yfiber

import (
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	onceSpan sync.Once
	span     trace.Tracer
)

func Span() trace.Tracer {
	onceSpan.Do(func() {
		span = otel.Tracer("fiber-server")
	})
	return span
}
