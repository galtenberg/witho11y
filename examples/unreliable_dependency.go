package main

import (
  "context"
  "fmt"

  "witho11y/pkg"
)

func ExampleUnreliableDependency(ctx context.Context, a string, b int) (int, string, error) {
  return 404, "You passed in: " + a, nil
}

func ObserveUnreliableDependency(ctx context.Context) {
  wrappedFunc := witho11y.BeforeAfterDurationWrapper(
    ExampleUnreliableDependency, witho11y.NewOTelTraceWrapper(), nil, nil)

  results, _ := wrappedFunc(ctx, "param1", 99)

  fmt.Println(results[1].(string))
}

func main() {
  ObserveUnreliableDependency(context.Background())
}
