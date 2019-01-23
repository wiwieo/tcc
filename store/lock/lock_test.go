package lock

import (
	"context"
	"testing"
	"time"
)

func TestEtcdLock_Lock(t *testing.T) {
	ctx := context.Background()
	l, err := NewEtcdLock(ctx, 6, "hello")

	if err != nil {
		t.Fatalf("it is failed to buy a locker through etcd, error info: %s", err)
	}
	println("it is success to buy a locker")

	err = l.Lock(ctx)
	if err != nil {
		t.Fatalf("it is failed to lock a locker, error info: %s", err)
	}
	println("it is success to lock a locker")
	time.Sleep(30 * time.Second)
	println("if locker is timeout")
	time.Sleep(30 * time.Second)
	err = l.Unlock(ctx)
	if err != nil {
		t.Fatalf("it is failed to unlock a locker, error info: %s", err)
	}
	println("it is success to unlock a locker")
}
