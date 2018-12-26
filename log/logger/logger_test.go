package logger

import "testing"

const path = `/mnt/d/project/go/src/log/file/log.log`

func TestNewStdLogger(t *testing.T) {
	l := NewStdLogger(true, true, true, true, true)
	l.SetPath(path)
	l.Trace("%s, %s", "hello", "world!")
}
