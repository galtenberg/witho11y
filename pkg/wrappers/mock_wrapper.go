package witho11y

import (
  "context"
  "sync"

  "witho11y/pkg"
)

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
