package otelmock

import (
  "fmt"

  "testing"
  "github.com/stretchr/testify/require"

  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
  traceSdk "go.opentelemetry.io/otel/sdk/trace"
)

func SetSpanAttributes(span trace.Span, params ...any) {
  attrs := make([]attribute.KeyValue, len(params))
  for i, param := range params {
    attrs[i] = attribute.String(fmt.Sprintf("param.%d", i), fmt.Sprintf("%v", param))
  }
  span.SetAttributes(attrs...)
}

func VerifySpanAttributes(t *testing.T, span traceSdk.ReadOnlySpan, expectedAttrs map[string]string) {
  attrs := span.Attributes()
  for k, v := range expectedAttrs {
    require.Contains(t, attrs, attribute.String(k,v))
  }
}
