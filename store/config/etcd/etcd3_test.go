package etcd3

import (
	"context"
	"testing"
	"time"
)

var st, err = NewEtcd3Client([]string{"localhost:2379"}, int(time.Minute), "", "", nil)

func TestNew(t *testing.T) {
	if err != nil{
		t.Fail()
	}
}

func TestEtcd3_Put(t *testing.T) {
	st.Put(context.Background(), "/tcc/api/222222", []byte("hello world"), 10)
}