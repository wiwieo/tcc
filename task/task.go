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
	if ts.off {
		ts.off = false
		go ts.exec()
	}
}

func (ts *Task) Stop() {
	ts.off = true
	ts.Off <- true
}

func (ts *Task) exec() {
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
	if len(data) == 0 {
		return
	}
	go taskToRetry(data)
	go taskToSend(data, "there is some exceptional data, please hurry up to resolve it")
}

// TODO 在使用levelDB时，因为数据没有共享，所以不存在并发问题
// 在使用共享数据（mysql）时， 分布式环境下，可能需要防止同时执行一个任务
func getBaseData() []*data.RequestInfo {
	// TODO 此处需要使用互斥锁 或者 简单起见， 只开一个任务
	needRollbackData, err := various.C.ListExceptionalRequestInfo()
	if err != nil {
		log.Errorf("the data that required for the task is failed to load, please check it. error information: %s", err)
		return nil
	}
	return needRollbackData
}
