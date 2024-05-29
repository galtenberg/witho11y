package witho11y

import (
  "context"
  "reflect"
)

func ConvertToReflectValues(params []any) []reflect.Value {
    values := make([]reflect.Value, len(params))
    for i, param := range params {
        values[i] = reflect.ValueOf(param)
    }
    return values
}

func CallWrapped(wrapped any, ctx context.Context, params []any) ([]reflect.Value, error) {
  wrValue := reflect.ValueOf(wrapped)
  wrParams := append([]reflect.Value{reflect.ValueOf(ctx)}, ConvertToReflectValues(params)...)
  results := wrValue.Call(wrParams)
  return results, nil
}

func ExtractResults(results []reflect.Value) ([]any, error) {
  finalResults := make([]any, len(results))
  var finalErr error
  for i, result := range results {
    if err, ok := result.Interface().(error); ok && err != nil {
      finalErr = err
    }
    finalResults[i] = result.Interface()
  }
  return finalResults, finalErr
}
