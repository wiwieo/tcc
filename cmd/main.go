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
	http.HandleFunc(fmt.Sprintf("%s/", serverName), tcc)

	go task.Start()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
