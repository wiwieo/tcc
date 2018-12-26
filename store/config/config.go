package config

import (
	"context"
)


type ConfigStore interface {
	Get(ctx context.Context, key string) ([]byte, error)
	List(ctx context.Context, prefix string) ([][]byte, error)
	WatchTree(ctx context.Context, prefix string, stopCh <-chan struct{}) ([]byte, error)

	Close(ctx context.Context)
}
