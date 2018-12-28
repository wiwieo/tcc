package task

import (
	"fmt"
	"strings"
	"tcc_transaction/constant"
	"tcc_transaction/global"
	"tcc_transaction/send"
	"tcc_transaction/send/email"
	"tcc_transaction/store/data"
)

func taskToSend(needRollbackData []*data.RequestInfo, subject string) {
	var s send.Send = email.NewEmailSender(*global.EmailUsername, subject, strings.Split(*global.EmailTo, ","))
	for _, v := range needRollbackData {
		if v.Times >= constant.RetryTimes && v.IsSend != constant.SendSuccess {
			err := s.Send([]byte(fmt.Sprintf("this data is wrong, please check it. information: %+v", v)))
			if err == nil {
				global.C.UpdateRequestInfoSend(v.Id)
			}
		}
	}
	if len(needRollbackData) > *global.MaxExceptionalData {
		s.Send([]byte(fmt.Sprintf("The exceptional data is too much [%d], please check it.", len(needRollbackData))))
	}
}
