package witho11y

import (
  "context"
  "testing"
  "fmt"
  "sync"

  util "witho11y/internal/util"
  "witho11y/pkg"

  "github.com/stretchr/testify/mock"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/assert"

  //"go.opentelemetry.io/otel/attribute"
)

//func Test_UnreliableDependency_WithTraceSpanOtel_Success(t *testing.T) {
  //sr, _ := util.SetupTestTrace()

  //mockBusinessLogic := &util.MockBusinessLogic{}
  //mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{"result1", "result2"}, nil)

  //wrappedLogic := witho11y.WithTraceSpanOtel("observe-reliable", mockBusinessLogic.Execute)
  //_, err := wrappedLogic(context.Background(), "param1", 99)
  //require.NoError(t, err)

  //spans := sr.Ended()
  //assert.Len(t, spans, 1)
  //require.Equal(t, "observe-reliable", spans[0].Name())
  //fmt.Println(spans[0])
  //util.VerifySpanAttributes(t, spans[0], map[string]string{ "param.0": "param1", "param.1": "99" })

  //mockBusinessLogic.AssertExpectations(t)
//}

//func Test_UnreliableDependency_WithTraceSpanOtel_Error(t *testing.T) {
  //sr, _ := util.SetupTestTrace()

  //mockBusinessLogic := &util.MockBusinessLogic{}
  //mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{nil}, fmt.Errorf("an error occurred"))

  //wrappedLogic := witho11y.WithTraceSpanOtel("observe-unreliable", mockBusinessLogic.Execute)
  //_, err := wrappedLogic(context.Background(), "param1", 99)

  //require.Error(t, err)

  //spans := sr.Ended()
  //assert.Len(t, spans, 1)
  //require.Equal(t, "observe-unreliable", spans[0].Name())
  //require.False(t, spans[0].EndTime().IsZero(), "expected span to be ended")
  //util.VerifySpanAttributes(t, spans[0], map[string]string{ "param.0": "param1", "param.1": "99" })

  //events := spans[0].Events()
  //require.Len(t, events, 1)
  //require.Equal(t, "exception", events[0].Name)
  //require.Contains(t, events[0].Attributes, attribute.String("exception.message", "an error occurred"))

  //mockBusinessLogic.AssertExpectations(t)
//}

type MockWrapper struct {
    events []witho11y.EventData
    mu     sync.Mutex
}

func NewMockWrapper() *MockWrapper {
    return &MockWrapper{
        events: []witho11y.EventData{},
    }
}

func (m *MockWrapper) Setup(ctx context.Context, name string) context.Context {
    event := witho11y.EventData{Name: name, Fields: make(map[string]interface{}), Ended: false}
    m.mu.Lock()
    m.events = append(m.events, event)
    m.mu.Unlock()
    ctx = context.WithValue(ctx, "eventName", name)
    return ctx
}

func (m *MockWrapper) AddFields(ctx context.Context, fields map[string]interface{}) {
    eventName := ctx.Value("eventName").(string)
    m.mu.Lock()
    defer m.mu.Unlock()
    for i := range m.events {
        if m.events[i].Name == eventName {
            for k, v := range fields {
                m.events[i].Fields[k] = v
            }
        }
    }
}

func (m *MockWrapper) Finish(ctx context.Context) {
    eventName := ctx.Value("eventName").(string)
    m.mu.Lock()
    defer m.mu.Unlock()
    for i := range m.events {
        if m.events[i].Name == eventName {
            m.events[i].Ended = true
        }
    }
}

func (m *MockWrapper) GetEvents() []witho11y.EventData {
    m.mu.Lock()
    defer m.mu.Unlock()
    return m.events
}

func Test_UnreliableDependency_WithTraceSpanOtel_Success(t *testing.T) {
  mockWrapper := NewMockWrapper()
  mockBusinessLogic := &util.MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{"result1", "result2"}, nil)

  beforeFields := map[string]interface{}{"param.0": nil, "param.1": nil}
  afterFields := map[string]interface{}{"post_param": "value2"}

  wrappedLogic := witho11y.WithTraceSpanOtel(mockBusinessLogic.Execute, mockWrapper, beforeFields, afterFields)
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
  mockWrapper := NewMockWrapper()
  mockBusinessLogic := &util.MockBusinessLogic{}
  mockBusinessLogic.On("Execute", mock.Anything, mock.Anything).Return([]any{nil}, fmt.Errorf("an error occurred"))

  beforeFields := map[string]interface{}{"param.0": nil, "param.1": nil}
  afterFields := map[string]interface{}{"post_param": "value2"}

  wrappedLogic := witho11y.WithTraceSpanOtel(mockBusinessLogic.Execute, mockWrapper, beforeFields, afterFields)
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
