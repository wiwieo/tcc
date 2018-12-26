package task

import (
	"tcc_transaction/log"
	"tcc_transaction/store/data"
	"tcc_transaction/store/data/mysql"
	"time"
)

var (
	c = mysql.NewMysqlClient("tcc", "tcc_123", "localhost", "3306", "tcc")
)

func TimerToExcuteTask() {
	t := time.NewTicker(time.Minute)
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
	needRollbackData, err := c.ListExceptionalRequestInfo()
	if err != nil {
		log.Errorf("the data that required for the task is failed to load, please check it. error information: %s", err)
		return nil
	}
	return needRollbackData
}
