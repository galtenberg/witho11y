package otelmock

import (
  "context"
  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
)

func TryUnreliableDependency(ctx context.Context, dep UnreliableDependency) error {
  tracer := otel.Tracer("example-tracer")
  var span trace.Span
  ctx, span = tracer.Start(ctx, "ExecuteOperation")
  defer span.End()

  result, err := dep.CallUnreliableDependency(ctx)
  if err != nil {
    span.SetAttributes(attribute.String("dependency.status", "failed"))
    span.RecordError(err)
    return err
  }

  span.SetAttributes(attribute.String("dependency.status", "succeeded"), attribute.String("dependency.result", result))
  return nil
}
