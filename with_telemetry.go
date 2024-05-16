package otelmock

import (
  "context"
  "fmt"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
)

func ObserveUnreliableDependency(ctx context.Context, dep UnreliableDependency, params ...interface{}) error {
  tracer := otel.Tracer("example-tracer")
  var span trace.Span
  ctx, span = tracer.Start(ctx, "ExecuteOperation")
  defer span.End()

  result, err := dep.CallUnreliableDependency(ctx, params...)
  if err != nil {
    span.SetAttributes(attribute.String("dependency.status", "failed"))
    span.RecordError(err)
    return err
  }

  span.SetAttributes(attribute.String("dependency.status", "succeeded"), attribute.String("dependency.result", result))
  return nil
}


func WithTelemetry(spanName string, businessLogic func(ctx context.Context, params ...interface{}) error) func(ctx context.Context, params ...interface{}) error {
  return func(ctx context.Context, params ...interface{}) error {
    tracer := otel.Tracer("example-tracer")
    ctx, span := tracer.Start(ctx, spanName)
    defer span.End()

    attrs := make([]attribute.KeyValue, len(params))
    for i, param := range params {
      attrs[i] = attribute.String(fmt.Sprintf("param.%d", i), fmt.Sprintf("%v", param))
    }
    span.SetAttributes(attrs...)

    err := businessLogic(ctx, params...)
    if err != nil {
      span.SetAttributes(attribute.String("dependency.status", "failed"))
      span.RecordError(err)
    }

    span.SetAttributes(attribute.String("dependency.status", "succeeded"))
    //span.SetAttributes(attribute.String("dependency.status", "succeeded"), attribute.String("dependency.result", result))
    return err
  }
}

func ExampleBusinessLogic(ctx context.Context, params ...interface{}) error {
  return nil
}

func ObserveUnreliableDependency2() {
  wrappedLogic := WithTelemetry("example-span", ExampleBusinessLogic)
  //err := wrappedLogic(context.Background(), "param1", 42)
  wrappedLogic(context.Background(), "param1", 42)
}
