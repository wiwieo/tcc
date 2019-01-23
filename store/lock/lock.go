package lock

import (
	"context"
)

type Lock interface {
	Lock(ctx context.Context) error
	Value(ctx context.Context, key string) ([]byte, error)
	Unlock(ctx context.Context) error
}
