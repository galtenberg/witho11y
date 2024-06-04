package witho11y

import (
  "context"
  "fmt"
)

type PrintScreen struct {
  events []EventData
}

func NewPrintScreen() *PrintScreen {
  return &PrintScreen{
    events: []EventData{},
  }
}

func (p *PrintScreen) Setup(ctx context.Context, name string) context.Context {
  event := EventData{Name: name, Fields: make(map[string]interface{}), Ended: false}
  p.events = append(p.events, event)
  fmt.Printf("Setup: %s\n", name)
  ctx = context.WithValue(ctx, "eventName", name)
  return ctx
}

func (p *PrintScreen) AddFields(ctx context.Context, fields map[string]interface{}) {
  eventName := ctx.Value("eventName").(string)
  for i := range p.events {
    if p.events[i].Name == eventName {
      for k, v := range fields {
        p.events[i].Fields[k] = v
      }
    }
  }
  fmt.Printf("AddFields: %v\n", fields)
}

func (p *PrintScreen) Finish(ctx context.Context) {
  eventName := ctx.Value("eventName").(string)
  for i := range p.events {
    if p.events[i].Name == eventName {
      p.events[i].Ended = true
    }
  }
  fmt.Printf("Finish: %s\n", eventName)
}

func (p *PrintScreen) GetEvents() []EventData {
  return p.events
}
