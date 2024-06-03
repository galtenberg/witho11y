package witho11y

import (
  "context"
  "time"

  "fmt"
  "reflect"
  "runtime"

  util "witho11y/internal/util"
  "witho11y/pkg"
)

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

func BeforeAfterDurationWrapper(wrappedFunc any, telemetry witho11y.TelemetryEvents, beforeFields, afterFields map[string]interface{}) func(ctx context.Context, params ...any) ([]any, error) {
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

    results, _ := util.CallWrapped(wrappedFunc, ctx, params)
    duration := time.Since(startTime)

    telemetry.AddFields(ctx, filterFields(fields, afterFields))

    telemetry.AddFields(ctx, map[string]interface{}{
      "dependency.status": "succeeded",
      "duration_ms":       float64(duration.Milliseconds()),
    })

    ret, err := util.ExtractResults(results)
    if err != nil {
      telemetry.AddFields(ctx, map[string]interface{}{
        "error": err.Error(),
      })
    }

    return ret, err
  }
}
