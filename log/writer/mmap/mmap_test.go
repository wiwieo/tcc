// +build linux,cgo darwin,cgo

package mmap

import (
	"fmt"
	"testing"
	"time"
)

const path = `/mnt/d/project/go/src/log/file/log.log`

func TestMmapWrite(t *testing.T) {
	m, err := NewMmap(path, 1<<14)
	if err != nil {
		panic(fmt.Sprintf("memory mapping to file error. %s", err))
	}
	for i := 0; i < 100000; i++ {
		err = m.Write([]byte(fmt.Sprintf("I haven't seen you for ages, %d.\r\n", i)))
		if err != nil {
			println(fmt.Sprintf("write to file failed. %s", err))
		}
	}
	time.Sleep(1 * time.Second)
	m.Close()
}

func BenchmarkMmapWrite(b *testing.B) {
	m, err := NewMmap(path, 1<<14)
	if err != nil {
		panic(fmt.Sprintf("memory mapping to file error. %s", err))
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err = m.Write([]byte("I haven't seen you for ages.\r\n"))
			if err != nil {
				println(fmt.Sprintf("write to file failed. %s", err))
				b.Fail()
			}
		}
	})
	m.Close()
}

func TestTime(t *testing.T) {
	now := time.Now()
	dest := time.Date(now.Year(), now.Month(), now.Day()+1, 1, 0, 0, 0, time.Local)
	println(fmt.Sprintf("%+v", dest))
}
