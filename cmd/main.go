package main

import (
  "context"
  "fmt"

  "otelmock/pkg"
)

func ExampleUnreliableDependency(ctx context.Context, params ...any) (int, string, error) {
  return 404, "try again", nil
}

func ObserveUnreliableDependency() {
  wrappedFunc := otelmock.WithTraceSpan("observe-unreliable-1", ExampleUnreliableDependency)
  results, _ := wrappedFunc(context.Background(), "param1", 42)
  fmt.Println(results)
}

func main() {
  ObserveUnreliableDependency()
}
