package main

import (
  "context"
  "fmt"

  telemeters "witho11y/pkg/telemeters"
  wrappers "witho11y/pkg/wrappers"
)

func ExampleUnreliableDependency(ctx context.Context, a string, b int) (int, string, error) {
  return 404, "You passed in: " + a, nil
}

func ObserveUnreliableDependency(ctx context.Context) {
  wrappedFunc := wrappers.BeforeAfterDurationWrapper(
    ExampleUnreliableDependency, telemeters.NewPrintScreen(), nil, nil)

  results, _ := wrappedFunc(ctx, "param1", 99)

  fmt.Println(results[1].(string))
}

func main() {
  ObserveUnreliableDependency(context.Background())
}
