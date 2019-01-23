package lock

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3/concurrency"
	"tcc_transaction/global/various"
	"time"
)

type etcdLock struct {
	mux    *concurrency.Mutex
	ttl    int
	prefix string
}

func NewEtcdLock(ctx context.Context, ttl int, prefix string) (Lock, error) {
	s, err := concurrency.NewSession(various.EtcdC.C, concurrency.WithTTL(ttl))
	if err != nil {
		return nil, err
	}
	//s.Orphan()
	l := concurrency.NewMutex(s, prefix)
	return &etcdLock{mux: l, ttl: ttl, prefix: prefix}, nil
}

func (e *etcdLock) Lock(ctx context.Context) error {
	// 是否需要做 超时获取锁 失败操作
	ctx, _ = context.WithTimeout(ctx, time.Duration(e.ttl)*time.Second)
	err := e.mux.Lock(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (e *etcdLock) Unlock(ctx context.Context) error {
	return e.mux.Unlock(ctx)
}

func (e *etcdLock) Value(ctx context.Context, key string) ([]byte, error) {
	return nil, fmt.Errorf("not support to retrive value of a lock")
}
