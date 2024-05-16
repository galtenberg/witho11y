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

// MockBusinessLogic is a mock implementation of the business logic function.
type MockBusinessLogic struct {
  mock.Mock
}

func (m *MockBusinessLogic) Execute(ctx context.Context, params ...interface{}) error {
  args := m.Called(ctx, params)
  return args.Error(0)
}

func TestWithTelemetry_Success(t *testing.T) {
  sr := tracetest.NewSpanRecorder()
  tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
  otel.SetTracerProvider(tp)

  mockBusinessLogic := &MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return(nil)

  wrappedLogic := WithTelemetry("example-span", mockBusinessLogic.Execute)
  err := wrappedLogic(context.Background(), "param1", 42)

  require.NoError(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)

  span := spans[0]
  attrs := span.Attributes()
  require.Contains(t, attrs, attribute.String("dependency.status", "succeeded"))

  require.Equal(t, "example-span", span.Name())
  require.Contains(t, attrs, attribute.String("param.0", "param1"))
  require.Contains(t, attrs, attribute.String("param.1", "42"))

  mockBusinessLogic.AssertExpectations(t)
}

func TestWithTelemetry_Error(t *testing.T) {
  sr := tracetest.NewSpanRecorder()
  tp := trace.NewTracerProvider(trace.WithSpanProcessor(sr))
  otel.SetTracerProvider(tp)

  mockBusinessLogic := &MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return(fmt.Errorf("an error occurred"))

  wrappedLogic := WithTelemetry("example-span", mockBusinessLogic.Execute)
  err := wrappedLogic(context.Background(), "param1", 42)

  require.Error(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)
  span := spans[0]

  require.Equal(t, "example-span", span.Name())
  require.False(t, span.EndTime().IsZero(), "expected span to be ended")

  attrs := span.Attributes()
  require.Contains(t, attrs, attribute.String("param.0", "param1"))
  require.Contains(t, attrs, attribute.String("param.1", "42"))

  events := span.Events()
  require.Len(t, events, 1)
  event := events[0]
  require.Equal(t, "exception", event.Name)
  require.Contains(t, event.Attributes, attribute.String("exception.message", "an error occurred"))

  mockBusinessLogic.AssertExpectations(t)
}
