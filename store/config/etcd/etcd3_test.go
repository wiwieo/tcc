package etcd3

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var st, err = NewEtcd3Client([]string{"localhost:2379"}, int(time.Minute), "", "", nil)

func TestNew(t *testing.T) {
	if err != nil {
		t.Fail()
	}
}

func TestEtcd3_Put(t *testing.T) {
	st.Put(context.Background(), "/tcc/api/222222", []byte("hello world"), 10)
	time.Sleep(time.Minute)
}

func TestXxx(t *testing.T) {
	ary := []int{0, 1, 2, 3, 4, 5}
	ary = append(ary[:1], ary[2:]...)
	println(fmt.Sprintf("%+v", ary))
}
