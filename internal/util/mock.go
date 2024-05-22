package otelmock

import (
  "context"

  "github.com/stretchr/testify/mock"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/sdk/trace"
  "go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type MockBusinessLogic struct {
  mock.Mock
}

func (m *MockBusinessLogic) Execute(ctx context.Context, params ...any) ([]any, error) {
  args := m.Called(ctx, params)
  return args.Get(0).([]any), args.Error(1)
}

func SetupTrace() (*tracetest.SpanRecorder, *trace.TracerProvider) {
  sr := tracetest.NewSpanRecorder()
  tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
  otel.SetTracerProvider(tp)
  return sr, tp
}
