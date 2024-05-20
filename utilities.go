package otelmock

import (
  "context"
  "fmt"
  "reflect"

  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
)

func setSpanAttributes(span trace.Span, params ...any) {
    attrs := make([]attribute.KeyValue, len(params))
    for i, param := range params {
        attrs[i] = attribute.String(fmt.Sprintf("param.%d", i), fmt.Sprintf("%v", param))
    }
    span.SetAttributes(attrs...)
}

func convertToReflectValues(params []any) []reflect.Value {
    values := make([]reflect.Value, len(params))
    for i, param := range params {
        values[i] = reflect.ValueOf(param)
    }
    return values
}

func callWrapped(wrapped any, ctx context.Context, params []any) ([]reflect.Value, error) {
  wrValue := reflect.ValueOf(wrapped)
  wrParams := append([]reflect.Value{reflect.ValueOf(ctx)}, convertToReflectValues(params)...)
  results := wrValue.Call(wrParams)
  return results, nil
}

func extractResults(results []reflect.Value, span trace.Span) ([]any, error) {
  finalResults := make([]any, len(results))
  var finalErr error
  for i, result := range results {
    if err, ok := result.Interface().(error); ok && err != nil {
      finalErr = err
      span.RecordError(err)
    }
  finalResults[i] = result.Interface()
  }
  return finalResults, finalErr
}
