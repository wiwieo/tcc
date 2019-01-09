package mysql

import (
	"fmt"
	"regexp"
	"testing"
)

func TestNewMysqlClient(t *testing.T) {
	db, _ := NewMysqlClient("tcc", "tcc_123", "localhost", "3306", "tcc")
	des, err := db.ListExceptionalRequestInfo()
	println(fmt.Sprintf("%+v, %s", des, err))
}

func TestReg(t *testing.T) {
	reg, err := regexp.Compile("^accounts/order/(.)*")
	if err == nil {
		println(reg.MatchString("accounts/order/1"))
	}
}
