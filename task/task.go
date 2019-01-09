package task

import (
	"tcc_transaction/global/config"
	"tcc_transaction/global/various"
	"tcc_transaction/log"
	"tcc_transaction/store/data"
	"time"
)

type Task struct {
	Interval time.Duration // unit: second
	Off      chan bool     // 控制停止任务
	off      bool          // 标记当前任务是否已经停止
	F        func()
}

func (ts *Task) Start() {
	if !ts.off {
		go ts.Exec()
	}
}

func (ts *Task) Stop() {
	ts.off = true
	ts.Off <- true
}

func (ts *Task) Exec() {
	t := time.NewTicker(time.Second * ts.Interval)
FOR:
	for {
		select {
		case <-t.C:
			go ts.F()
		case off := <-ts.Off:
			if off {
				break FOR
			}
		}
	}
}

var defaultTask = &Task{
	Interval: time.Duration(*config.TimerInterval),
	Off:      make(chan bool, 1),
	F:        retryAndSend,
}

func Start() {
	defaultTask.Start()
}

func Stop() {
	defaultTask.Stop()
}

func retryAndSend() {
	data := getBaseData()
	go taskToRetry(data)
	go taskToSend(data, "there is some exceptional data, please hurry up to resolve it")
}

func getBaseData() []*data.RequestInfo {
	needRollbackData, err := various.C.ListExceptionalRequestInfo()
	if err != nil {
		log.Errorf("the data that required for the task is failed to load, please check it. error information: %s", err)
		return nil
	}
	return needRollbackData
}
