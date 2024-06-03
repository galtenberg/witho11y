package witho11y

import (
  "context"
  "time"

    "fmt"
    "reflect"
    "runtime"
    "sync"

  "witho11y/internal/util"

  "go.opentelemetry.io/otel"
  "go.opentelemetry.io/otel/trace"
  "go.opentelemetry.io/otel/attribute"
)

type TelemetryEvents interface {
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

type OTelTraceWrapper struct {
  tracer trace.Tracer
  events []EventData
  mu     sync.Mutex
}

func NewOTelTraceWrapper() *OTelTraceWrapper {
  return &OTelTraceWrapper{
    tracer: otel.Tracer("example-tracer"),
    events: []EventData{},
  }
}

func (o *OTelTraceWrapper) Setup(ctx context.Context, name string) context.Context {
  var span trace.Span
  ctx, span = o.tracer.Start(ctx, name)
  o.mu.Lock()
  o.events = append(o.events, EventData{Name: name, Fields: make(map[string]interface{}), Ended: false})
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

func (o *OTelTraceWrapper) GetEvents() []EventData {
  o.mu.Lock()
  defer o.mu.Unlock()
  return o.events
}

func filterFields(fields, subset map[string]interface{}) map[string]interface{} {
  if subset == nil {
    return fields
  }
  if len(subset) == 0 {
    return nil
  }
  filtered := make(map[string]interface{})
  for k, v := range subset {
    if val, ok := fields[k]; ok {
      filtered[k] = val
    } else if v != nil {
      filtered[k] = v
    }
  }
  return filtered
}

func BeforeAfterDurationWrapper(wrappedFunc any, telemetry TelemetryEvents, beforeFields, afterFields map[string]interface{}) func(ctx context.Context, params ...any) ([]any, error) {
  return func(ctx context.Context, params ...any) ([]any, error) {
    funcName := runtime.FuncForPC(reflect.ValueOf(wrappedFunc).Pointer()).Name()
    spanName := fmt.Sprintf("%s-%d", funcName, time.Now().UnixNano())

    ctx = telemetry.Setup(ctx, spanName)
    defer telemetry.Finish(ctx)

    startTime := time.Now()

    fields := make(map[string]interface{})
    for i, param := range params {
      fields[fmt.Sprintf("param.%d", i)] = param
    }

    telemetry.AddFields(ctx, filterFields(fields, beforeFields))

    results, _ := witho11y.CallWrapped(wrappedFunc, ctx, params)
    duration := time.Since(startTime)

    telemetry.AddFields(ctx, filterFields(fields, afterFields))

    telemetry.AddFields(ctx, map[string]interface{}{
      "dependency.status": "succeeded",
      "duration_ms":       float64(duration.Milliseconds()),
    })

    ret, err := witho11y.ExtractResults(results)
    if err != nil {
      telemetry.AddFields(ctx, map[string]interface{}{
        "error": err.Error(),
      })
    }

    return ret, err
  }
}
