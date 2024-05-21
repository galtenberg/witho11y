package otelmock

import (
  "context"
  "time"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
)

func WithTelemetry(spanName string, wrappedFunc any) func(ctx context.Context, params ...any) ([]any, error) {
  return func(ctx context.Context, params ...any) ([]any, error) {
    tracer := otel.Tracer("observe-tracer")
    ctx, span := tracer.Start(ctx, spanName)
    defer span.End()

    startTime := time.Now()
    setSpanAttributes(span, params...)

    results, err := callWrapped(wrappedFunc, ctx, params)
    duration := time.Since(startTime)
    if err != nil {
      return nil, err
    }

    span.SetAttributes(attribute.String("dependency.status", "succeeded"), attribute.Float64("duration_ms", float64(duration.Milliseconds())))
    return extractResults(results, span)
  }
}
