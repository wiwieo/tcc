package main

import (
	"fmt"
	"net/http"
	"tcc_transaction/global"
	"tcc_transaction/task"
)

var serverName = "/tcc"

func main() {
	global.InitAll()

	http.Handle("/", http.FileServer(http.Dir("file")))
	http.HandleFunc(fmt.Sprintf("%s/", serverName), tcc)

	go task.TimerToExcuteTask()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
