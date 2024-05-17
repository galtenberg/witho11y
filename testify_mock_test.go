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

  span := spans[0]
  require.Equal(t, "observe-reliable", span.Name())

  expectedAttrs := map[string]string {
    "param.0": "param1",
    "param.1": "42",
  }
  verifySpanAttributes(t, span, expectedAttrs)

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
  span := spans[0]

  require.Equal(t, "observe-unreliable", span.Name())
  require.False(t, span.EndTime().IsZero(), "expected span to be ended")

  expectedAttrs := map[string]string {
    "param.0": "param1",
    "param.1": "42",
  }
  verifySpanAttributes(t, span, expectedAttrs)

  events := span.Events()
  require.Len(t, events, 1)
  event := events[0]
  require.Equal(t, "exception", event.Name)
  require.Contains(t, event.Attributes, attribute.String("exception.message", "an error occurred"))

  mockBusinessLogic.AssertExpectations(t)
}
