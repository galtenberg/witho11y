package witho11y

import (
  "context"

  "fmt"
  "sync"

  "witho11y/pkg"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/attribute"
)

type OTelTraceWrapper struct {
  tracer trace.Tracer
  events []witho11y.EventData
  mu     sync.Mutex
}

func NewOTelTraceWrapper() *OTelTraceWrapper {
  return &OTelTraceWrapper{
    tracer: otel.Tracer("example-tracer"),
    events: []witho11y.EventData{},
  }
}

func (o *OTelTraceWrapper) Setup(ctx context.Context, name string) context.Context {
  var span trace.Span
  ctx, span = o.tracer.Start(ctx, name)
  o.mu.Lock()
  o.events = append(o.events, witho11y.EventData{Name: name, Fields: make(map[string]interface{}), Ended: false})
  o.mu.Unlock()
  ctx = context.WithValue(ctx, "spanName", name)
  return context.WithValue(ctx, "span", span)
}

func (o *OTelTraceWrapper) AddFields(ctx context.Context, fields map[string]interface{}) {
  span := trace.SpanFromContext(ctx)
  attrs := make([]attribute.KeyValue, 0, len(fields))
  for k, v := range fields {
    attrs = append(attrs, attribute.String(k, fmt.Sprintf("%v", v)))
  }
  span.SetAttributes(attrs...)

  spanName := ctx.Value("spanName").(string)
  o.mu.Lock()
  defer o.mu.Unlock()
  for i := range o.events {
    if o.events[i].Name == spanName {
      for k, v := range fields {
        o.events[i].Fields[k] = v
      }
    }
  }
}

func (o *OTelTraceWrapper) Finish(ctx context.Context) {
  span := trace.SpanFromContext(ctx)
  span.End()

  spanName := ctx.Value("spanName").(string)
  o.mu.Lock()
  defer o.mu.Unlock()
  for i := range o.events {
    if o.events[i].Name == spanName {
      o.events[i].Ended = true
    }
  }
}

func (o *OTelTraceWrapper) GetEvents() []witho11y.EventData {
  o.mu.Lock()
  defer o.mu.Unlock()
  return o.events
}
