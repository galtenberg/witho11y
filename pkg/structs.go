package witho11y

import (
  "context"
)

type Telemeter interface {
  Setup(ctx context.Context, name string) context.Context
  AddFields(ctx context.Context, fields map[string]interface{})
  Finish(ctx context.Context)
  GetEvents() []EventData
}

type EventData struct {
  Name   string
  Fields map[string]interface{}
  Ended  bool
}
