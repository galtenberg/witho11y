package otelmock

import (
  "context"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
)

func WithTelemetry(spanName string, businessLogic any) func(ctx context.Context, params ...any) ([]any, error) {
  return func(ctx context.Context, params ...any) ([]any, error) {
    tracer := otel.Tracer("observe-tracer")
    ctx, span := tracer.Start(ctx, spanName)
    defer span.End()

    setSpanAttributes(span, params...)

    results, err := callWrapped(businessLogic, ctx, params)
    if err != nil {
      return nil, err
    }

    span.SetAttributes(attribute.String("dependency.status", "succeeded"))
    return extractResults(results, span)
  }
}

func ExampleBusinessLogic(ctx context.Context, params ...interface{}) error {
  return nil
}

func ObserveUnreliableDependency2() {
  wrappedLogic := WithTelemetry("observe-unreliable-1", ExampleBusinessLogic)
  //err := WithTelemetry(context.Background(), "param1", 42)
  wrappedLogic(context.Background(), "param1", 42)
}
