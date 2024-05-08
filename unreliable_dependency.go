package otelmock

import "context"

type UnreliableDependency interface {
  CallUnreliableDependency(ctx context.Context, params ...interface{}) (string, error)
}
