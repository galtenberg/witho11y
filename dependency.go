package otelmock

import "context"

type Dependency interface {
  CallDependency(ctx context.Context) (string, error)
}
