package otelmock

import (
  "context"
  "testing"
  "fmt"

  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/assert"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/sdk/trace"
  "go.opentelemetry.io/otel/sdk/trace/tracetest"
)

type MockBusinessLogic struct {
  mock.Mock
}

func (m *MockBusinessLogic) Execute(ctx context.Context, params ...interface{}) error {
  args := m.Called(ctx, params)
  return args.Error(0)
}

func setupTrace() (*tracetest.SpanRecorder, *trace.TracerProvider) {
  sr := tracetest.NewSpanRecorder()
  tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
  otel.SetTracerProvider(tp)
  return sr, tp
}

func verifySpanAttributes(t *testing.T, span trace.ReadOnlySpan, expectedAttrs map[string]string) {
  attrs := span.Attributes()
  for k, v := range expectedAttrs {
    require.Contains(t, attrs, attribute.String(k,v))
  }
}

func TestWithTelemetry_Success(t *testing.T) {
  sr, _ := setupTrace()

  mockBusinessLogic := &MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return(nil)

  wrappedLogic := WithTelemetry("observe-reliable", mockBusinessLogic.Execute)
  err := wrappedLogic(context.Background(), "param1", 42)

  require.NoError(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)
  require.Equal(t, "observe-reliable", spans[0].Name())
  verifySpanAttributes(t, spans[0], map[string]string{ "param.0": "param1", "param.1": "42" })

  mockBusinessLogic.AssertExpectations(t)
}

func TestWithTelemetry_Error(t *testing.T) {
  sr, _ := setupTrace()

  mockBusinessLogic := &MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return(fmt.Errorf("an error occurred"))

  wrappedLogic := WithTelemetry("observe-unreliable", mockBusinessLogic.Execute)
  err := wrappedLogic(context.Background(), "param1", 42)

  require.Error(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)
  require.Equal(t, "observe-unreliable", spans[0].Name())
  require.False(t, spans[0].EndTime().IsZero(), "expected span to be ended")
  verifySpanAttributes(t, spans[0], map[string]string{ "param.0": "param1", "param.1": "42" })

  events := spans[0].Events()
  require.Len(t, events, 1)
  require.Equal(t, "exception", events[0].Name)
  require.Contains(t, events[0].Attributes, attribute.String("exception.message", "an error occurred"))

  mockBusinessLogic.AssertExpectations(t)
}
