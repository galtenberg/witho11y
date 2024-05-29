package witho11y

import (
  "context"
  "time"

  "witho11y/internal/util"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
)

func WithTraceSpanOtel(spanName string, wrappedFunc any) func(ctx context.Context, params ...any) ([]any, error) {
  return func(ctx context.Context, params ...any) ([]any, error) {
    tracer := otel.Tracer("observe-tracer")
    ctx, span := tracer.Start(ctx, spanName)
    defer span.End()

    startTime := time.Now()
    witho11y.SetSpanAttributes(span, params...)

    results, _ := witho11y.CallWrapped(wrappedFunc, ctx, params)
    duration := time.Since(startTime)

    span.SetAttributes(
      attribute.String("dependency.status", "succeeded"),
      attribute.Float64("duration_ms", float64(duration.Milliseconds())),
    )

    ret, err := witho11y.ExtractResults(results)
    if err != nil {
      span.RecordError(err)
    }

    return ret, err
  }
}
