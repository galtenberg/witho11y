package otelmock

import (
  "context"
  "testing"
  "errors"

  "github.com/stretchr/testify/assert"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/sdk/trace"
  "go.opentelemetry.io/otel/sdk/trace/tracetest"

  gomock "go.uber.org/mock/gomock"
)

func AttributesToMap(attrs []attribute.KeyValue) map[attribute.Key]attribute.Value {
  attrMap := make(map[attribute.Key]attribute.Value)
  for _, attr := range attrs { attrMap[attr.Key] = attr.Value }
  return attrMap
}

func TestExecuteOperation_Success(t *testing.T) {
  ctrl := gomock.NewController(t)
  mockDep := NewMockDependency(ctrl)
  mockDep.EXPECT().CallDependency(gomock.Any()).Return("success result", nil)

  sr := tracetest.NewSpanRecorder()
  tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
  otel.SetTracerProvider(tp)

  ctx := context.Background()
  err := ExecuteOperation(ctx, mockDep)
  assert.NoError(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)
  assert.Equal(t, "succeeded", AttributesToMap(spans[0].Attributes())["dependency.status"].AsString())

  mockDep.EXPECT().CallDependency(gomock.Any()).Return("", errors.New("failure"))
  err = ExecuteOperation(ctx, mockDep)
  assert.Error(t, err)
}
