package main

import (
  "context"
  "fmt"

  "witho11y/pkg"
)

func ExampleUnreliableDependency(ctx context.Context, a string, b int) (int, string, error) {
  return 404, "You passed in: " + a, nil
}

func ObserveUnreliableDependency() {
  results, _ := witho11y.WithTraceSpanOtel("observe-unreliable-1",
    ExampleUnreliableDependency)(context.Background(), "param1", 99)
  //wrappedFunc := witho11y.WithTraceSpanOtel("observe-unreliable-1", ExampleUnreliableDependency)
  //results, _ := wrappedFunc(context.Background(), "param1", 99)

  fmt.Println(results[1].(string))
}

func main() {
  ObserveUnreliableDependency()
}
