package otelmock

import "context"

type UnreliableDependency interface {
  CallUnreliableDependency(ctx context.Context) (string, error)
}
