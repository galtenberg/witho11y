package witho11y

import (
  "context"
  "testing"
  "fmt"

  util "witho11y/internal/util"
  "witho11y/pkg"

  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/assert"

  "go.opentelemetry.io/otel/attribute"
)

func Test_UnreliableDependency_WithTraceSpanOtel_Success(t *testing.T) {
  sr, _ := util.SetupTestTrace()

  mockBusinessLogic := &util.MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{"result1", "result2"}, nil)

  wrappedLogic := witho11y.WithTraceSpanOtel("observe-reliable", mockBusinessLogic.Execute)
  _, err := wrappedLogic(context.Background(), "param1", 99)
  require.NoError(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)
  require.Equal(t, "observe-reliable", spans[0].Name())
  fmt.Println(spans[0])
  util.VerifySpanAttributes(t, spans[0], map[string]string{ "param.0": "param1", "param.1": "99" })

  mockBusinessLogic.AssertExpectations(t)
}

func Test_UnreliableDependency_WithTraceSpanOtel_Error(t *testing.T) {
  sr, _ := util.SetupTestTrace()

  mockBusinessLogic := &util.MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{nil}, fmt.Errorf("an error occurred"))

  wrappedLogic := witho11y.WithTraceSpanOtel("observe-unreliable", mockBusinessLogic.Execute)
  _, err := wrappedLogic(context.Background(), "param1", 99)

  require.Error(t, err)

  spans := sr.Ended()
  assert.Len(t, spans, 1)
  require.Equal(t, "observe-unreliable", spans[0].Name())
  require.False(t, spans[0].EndTime().IsZero(), "expected span to be ended")
  util.VerifySpanAttributes(t, spans[0], map[string]string{ "param.0": "param1", "param.1": "99" })

  events := spans[0].Events()
  require.Len(t, events, 1)
  require.Equal(t, "exception", events[0].Name)
  require.Contains(t, events[0].Attributes, attribute.String("exception.message", "an error occurred"))

  mockBusinessLogic.AssertExpectations(t)
}
