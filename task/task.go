package task

import (
	"tcc_transaction/global"
	"tcc_transaction/log"
	"tcc_transaction/store/data"
	"time"
)

func TimerToExcuteTask() {
	t := time.NewTicker(time.Second * time.Duration(*global.TimerInterval))
	for {
		select {
		case <-t.C:
			data := getBaseData()
			go taskToRetry(data)
			go taskToSend(data, "there is some exceptional data, please hurry up to resolve it")
		}
	}
}

func getBaseData() []*data.RequestInfo{
	needRollbackData, err := global.C.ListExceptionalRequestInfo()
	if err != nil {
		log.Errorf("the data that required for the task is failed to load, please check it. error information: %s", err)
		return nil
	}
	return needRollbackData
}
