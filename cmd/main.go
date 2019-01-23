package main

import (
	"fmt"
	"net/http"
	"tcc_transaction/global/various"
	"tcc_transaction/task"
)

var serverName = "/tcc"

func main() {
	various.InitAll()
	http.Handle("/", http.FileServer(http.Dir("file")))
	// 用于决定使用哪种tcc逻辑，自定义或默认
	var rtnHandle = func(t tcc) func(http.ResponseWriter, *http.Request) {
		p := &proxy{}
		p.t = t
		return p.process
	}
	http.HandleFunc("/tcc/examples/", rtnHandle(NewExampleTcc()))
	http.HandleFunc(fmt.Sprintf("%s/", serverName), rtnHandle(NewDefaultTcc()))

	go task.Start()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
