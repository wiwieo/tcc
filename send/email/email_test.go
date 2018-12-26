package email

import (
	"testing"
	"time"
)

func TestEmail_Send(t *testing.T) {
	e := NewEmailSender("qingwei.wu@uuabc.com", "program error info", []string{"qingwei.wu@uuabc.com"}, )
	e.Send([]byte("hello, your program is wrong, hurry up to execute, ape."))
	time.Sleep(time.Minute)
}