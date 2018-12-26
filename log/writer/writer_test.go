package writer

import (
	"fmt"
	"testing"
	"time"
)

const path = `/mnt/d/project/go/src/log/file/log.log`

func TestWrite(t *testing.T) {
	w := NewWriter(path, 1<<20)
	for i := 0; i < 1000000; i++ {
		err := w.Write([]byte(fmt.Sprintf("I haven't seen you for ages, %d.\r\n", i)))
		if err != nil {
			println(fmt.Sprintf("write to file failed. %s", err))
		}
	}
	time.Sleep(1 * time.Second)
	w.Close()
}
