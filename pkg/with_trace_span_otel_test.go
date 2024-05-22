package otelmock

import (
  "context"
  "testing"
  "fmt"

  util "otelmock/internal/util"
  "otelmock/pkg"

  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/assert"

  "go.opentelemetry.io/otel/attribute"
)

func TestWithTraceSpanOtel_Success(t *testing.T) {
  sr, _ := util.SetupTrace()

  mockBusinessLogic := &util.MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{"result1", "result2"}, nil)

  wrappedLogic := otelmock.WithTraceSpanOtel("observe-reliable", mockBusinessLogic.Execute)
  _, err := wrappedLogic(context.Background(), "param1", 42)
  require.NoError(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)
  require.Equal(t, "observe-reliable", spans[0].Name())
  fmt.Println(spans[0])
  util.VerifySpanAttributes(t, spans[0], map[string]string{ "param.0": "param1", "param.1": "42" })

  mockBusinessLogic.AssertExpectations(t)
}

func TestWithTraceSpanOtel_Error(t *testing.T) {
  sr, _ := util.SetupTrace()

  mockBusinessLogic := &util.MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{nil}, fmt.Errorf("an error occurred"))

  wrappedLogic := otelmock.WithTraceSpanOtel("observe-unreliable", mockBusinessLogic.Execute)
  _, err := wrappedLogic(context.Background(), "param1", 42)

  require.Error(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)
  require.Equal(t, "observe-unreliable", spans[0].Name())
  require.False(t, spans[0].EndTime().IsZero(), "expected span to be ended")
  util.VerifySpanAttributes(t, spans[0], map[string]string{ "param.0": "param1", "param.1": "42" })

  events := spans[0].Events()
  require.Len(t, events, 1)
  require.Equal(t, "exception", events[0].Name)
  require.Contains(t, events[0].Attributes, attribute.String("exception.message", "an error occurred"))

  mockBusinessLogic.AssertExpectations(t)
}

  //results, err := wrappedLogic(context.Background(), "param1", 42)
  //assert.Equal(t, []any{"result1", "result2"}, results[0])
  //result1, _ := results[0].(string)
  //assert.Equal(t, "result1", result1)
