package config

import (
	"context"
)

type ConfigStore interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Put(ctx context.Context, key string, value []byte, ttl int) ([]byte, error)
	List(ctx context.Context, prefix string) ([][]byte, error)
	WatchTree(ctx context.Context, prefix string, callback func(k, v []byte, tpe string))

	Close(ctx context.Context)
}
