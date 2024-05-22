package otelmock

import (
  "context"
  "time"

  "otelmock/internal/util"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
)

func WithTraceSpanOtel(spanName string, wrappedFunc any) func(ctx context.Context, params ...any) ([]any, error) {
  return func(ctx context.Context, params ...any) ([]any, error) {
    tracer := otel.Tracer("observe-tracer")
    ctx, span := tracer.Start(ctx, spanName)
    defer span.End()

    startTime := time.Now()
    otelmock.SetSpanAttributes(span, params...)

    results, err := otelmock.CallWrapped(wrappedFunc, ctx, params)
    duration := time.Since(startTime)
    if err != nil {
      span.RecordError(err)
      return nil, err
    }

    span.SetAttributes(attribute.String("dependency.status", "succeeded"), attribute.Float64("duration_ms", float64(duration.Milliseconds())))
    ret, finalErr := otelmock.ExtractResults(results)
    if finalErr != nil {
      span.RecordError(finalErr)
    }
    return ret, finalErr
  }
}
