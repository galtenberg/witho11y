package main

import (
  "context"
  "fmt"

  "otelmock/pkg"
)

func ExampleHighLatency(ctx context.Context, params ...any) (int, string, error) {
  return 404, "try again", nil
}

func ObserveHighLatency() {
  wrappedFunc := otelmock.WithTraceSpanOtel("observe-latency-1", ExampleHighLatency)
  results, _ := wrappedFunc(context.Background(), "param1", 99)
  fmt.Println(results[1].(string))
}

func main() {
  ObserveHighLatency()
}
