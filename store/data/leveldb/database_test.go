package leveldb

import (
	"fmt"
	"tcc_transaction/store/data"
	"testing"
)

func TestNewLevelDB(t *testing.T) {
	ldb, err := NewLevelDB("./tcc")
	if err != nil {
		panic(err)
	}
	ri := &data.RequestInfo{
		Method: "put",
		Url:    "/tcc/test/",
	}
	err = ldb.InsertRequestInfo(ri)
	if err != nil {
		panic(err)
	}
	rri, err := ldb.getRequestInfo(ri.Id)
	if err != nil {
		panic(err)
	}
	println(fmt.Sprintf("return request info from levelDB is : %+v", rri))
}

func TestLevelDBClient_ListExceptionalRequestInfo(t *testing.T) {
	ldb, err := NewLevelDB("./tcc")
	if err != nil {
		panic(err)
	}
	ris, err := ldb.ListExceptionalRequestInfo()
	if err != nil {
		panic(err)
	}
	println(len(ris))
}

func TestGenerateID(t *testing.T) {
	println(generateID())
}
