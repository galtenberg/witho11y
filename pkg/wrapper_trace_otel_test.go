package witho11y

import (
  "context"
  "testing"
  "fmt"

  util "witho11y/internal/util"
  wrappers "witho11y/pkg/wrappers"

  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/assert"
)

func Test_UnreliableDependency_WithTraceSpanOtel_Success(t *testing.T) {
  mockWrapper := wrappers.NewMockWrapper()
  mockBusinessLogic := &util.MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{"result1", "result2"}, nil)

  beforeFields := map[string]interface{}{"param.0": nil, "param.1": nil}
  afterFields := map[string]interface{}{"post_param": "value2"}

  wrappedLogic := wrappers.BeforeAfterDurationWrapper(mockBusinessLogic.Execute, mockWrapper, beforeFields, afterFields)
  _, err := wrappedLogic(context.Background(), "param1", 99)
  require.NoError(t, err)

  events := mockWrapper.GetEvents()
  assert.Len(t, events, 1)
  require.Contains(t, events[0].Name, "MockBusinessLogic")
  require.True(t, events[0].Ended)
  assert.Equal(t, map[string]interface{}{
    "param.0": "param1",
    "param.1": 99,
    "post_param": "value2",
    "dependency.status": "succeeded",
    "duration_ms": float64(0),
  }, events[0].Fields)

  mockBusinessLogic.AssertExpectations(t)
}

func Test_UnreliableDependency_WithTraceSpanOtel_Error(t *testing.T) {
  mockWrapper := wrappers.NewMockWrapper()
  mockBusinessLogic := &util.MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{nil}, fmt.Errorf("an error occurred"))

  beforeFields := map[string]interface{}{"param.0": nil, "param.1": nil}
  afterFields := map[string]interface{}{"post_param": "value2"}

  wrappedLogic := wrappers.BeforeAfterDurationWrapper(mockBusinessLogic.Execute, mockWrapper, beforeFields, afterFields)
  _, err := wrappedLogic(context.Background(), "param1", 99)
  require.Error(t, err)

  events := mockWrapper.GetEvents()
  assert.Len(t, events, 1)
  require.Contains(t, events[0].Name, "MockBusinessLogic")
  require.True(t, events[0].Ended)
  assert.Equal(t, map[string]interface{}{
    "param.0": "param1",
    "param.1": 99,
    "post_param": "value2",
    "dependency.status": "succeeded",
    "duration_ms": float64(0),
    "error": "an error occurred",
  }, events[0].Fields)

  mockBusinessLogic.AssertExpectations(t)
}
