package witho11y

import (
  "fmt"

  "testing"
  "github.com/stretchr/testify/require"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"
  sdktrace "go.opentelemetry.io/otel/sdk/trace"
  "go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func SetupTestTrace() (*tracetest.SpanRecorder, *sdktrace.TracerProvider) {
  sr := tracetest.NewSpanRecorder()
  tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
  otel.SetTracerProvider(tp)
  return sr, tp
}

func SetSpanAttributes(span trace.Span, params ...any) {
  attrs := make([]attribute.KeyValue, len(params))
  for i, param := range params {
    attrs[i] = attribute.String(fmt.Sprintf("param.%d", i), fmt.Sprintf("%v", param))
  }
  span.SetAttributes(attrs...)
}

func VerifySpanAttributes(t *testing.T, span sdktrace.ReadOnlySpan, expectedAttrs map[string]string) {
  attrs := span.Attributes()
  for k, v := range expectedAttrs {
    require.Contains(t, attrs, attribute.String(k,v))
  }
}
