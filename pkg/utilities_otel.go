package otelmock

import (
  "fmt"
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
