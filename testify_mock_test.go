package otelmock

import (
  "context"
  "testing"
  "fmt"

  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/attribute"
  "go.opentelemetry.io/otel/trace"

  //"go.opentelemetry.io/otel/trace/embedded"
)

// MockTracerProvider is a mock implementation of trace.TracerProvider
type MockTracerProvider struct {
  mock.Mock
}

//type MockTracerProvider interface{ tracerProvider() }

func (m *MockTracerProvider) Tracer(name string, opts ...trace.TracerOption) trace.Tracer {
  args := m.Called(name, opts)
  return args.Get(0).(trace.Tracer)
}

func (m *MockTracerProvider) tracerProvider() {}

// MockTracer is a mock implementation of trace.Tracer
type MockTracer struct {
  mock.Mock
}

func (m *MockTracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
  args := m.Called(ctx, name, opts)
  return args.Get(0).(context.Context), args.Get(1).(trace.Span)
}

// MockSpan is a mock implementation of trace.Span
type MockSpan struct {
  mock.Mock
}

func (m *MockSpan) End(options ...trace.SpanEndOption) {
  m.Called(options)
}

func (m *MockSpan) RecordError(err error, options ...trace.EventOption) {
  m.Called(err, options)
}

func (m *MockSpan) SetAttributes(attributes ...attribute.KeyValue) {
  m.Called(attributes)
}

// Other methods of trace.Span can be mocked as needed...

func TestWithTelemetry_Success(t *testing.T) {
  mockTracerProvider := &MockTracerProvider{}
  mockTracer := &MockTracer{}
  mockSpan := &MockSpan{}

  //var _ trace.TracerProvider = (*MockTracerProvider)(nil)
  otel.SetTracerProvider(mockTracerProvider)

  mockTracerProvider.On("Tracer", "example-tracer", mock.Anything).Return(mockTracer)
  mockTracer.On("Start", mock.Anything, "example-span", mock.Anything).Return(context.Background(), mockSpan)
  mockSpan.On("End", mock.Anything)
  mockSpan.On("SetAttributes", mock.Anything)
  mockSpan.On("RecordError", mock.Anything, mock.Anything)

  wrappedLogic := WithTelemetry("example-span", ExampleBusinessLogic)
  err := wrappedLogic(context.Background(), "param1", 42)

  require.NoError(t, err)
  mockTracerProvider.AssertExpectations(t)
  mockTracer.AssertExpectations(t)
  mockSpan.AssertExpectations(t)
}

func TestWithTelemetry_Error(t *testing.T) {
  mockTracerProvider := &MockTracerProvider{}
  mockTracer := &MockTracer{}
  mockSpan := &MockSpan{}

  otel.SetTracerProvider(mockTracerProvider)

  mockTracerProvider.On("Tracer", "example-tracer", mock.Anything).Return(mockTracer)
  mockTracer.On("Start", mock.Anything, "example-span", mock.Anything).Return(context.Background(), mockSpan)
  mockSpan.On("End", mock.Anything)
  mockSpan.On("SetAttributes", mock.Anything)
  mockSpan.On("RecordError", mock.Anything, mock.Anything)

  errorLogic := func(ctx context.Context, params ...interface{}) error {
    return fmt.Errorf("an error occurred")
  }

  wrappedLogic := WithTelemetry("example-span", errorLogic)
  err := wrappedLogic(context.Background(), "param1", 42)

  require.Error(t, err)
  mockTracerProvider.AssertExpectations(t)
  mockTracer.AssertExpectations(t)
  mockSpan.AssertExpectations(t)
}
