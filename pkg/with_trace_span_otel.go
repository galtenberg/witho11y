package otelmock

import (
  "context"
  "time"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
)

func WithTraceSpanOtel(spanName string, wrappedFunc any) func(ctx context.Context, params ...any) ([]any, error) {
  return func(ctx context.Context, params ...any) ([]any, error) {
    tracer := otel.Tracer("observe-tracer")
    ctx, span := tracer.Start(ctx, spanName)
    defer span.End()

    startTime := time.Now()
    setSpanAttributes(span, params...)

    results, err := callWrapped(wrappedFunc, ctx, params)
    duration := time.Since(startTime)
    if err != nil {
      span.RecordError(err)
      return nil, err
    }

    span.SetAttributes(attribute.String("dependency.status", "succeeded"), attribute.Float64("duration_ms", float64(duration.Milliseconds())))
    ret, finalErr := extractResults(results)
    if finalErr != nil {
      span.RecordError(finalErr)
    }
    return ret, finalErr
  }
}
