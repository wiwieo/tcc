package data

import "tcc_transaction/util"

type RequestInfo struct {
	Id           int64          `json:"id" db:"id"`
	Url          string         `json:"url" db:"url"`
	Method       string         `json:"method" db:"method"`
	Param        string         `json:"param" db:"param"`
	Status       int            `json:"status" db:"status"`
	Times        int            `json:"times" db:"times"`
	IsSend       int            `json:"is_send" db:"is_send"`
	Deleted      int            `json:"deleted" db:"deleted"`
	CreateTime   int64          `json:"create_time" db:"create_time"`
	UpdateTime   int64          `json:"update_time" db:"update_time"`
	SuccessSteps []*SuccessStep `json:"success_steps" `
}

type SuccessStep struct {
	Id         int64  `json:"id" db:"id"`
	RequestId  int64  `json:"request_id" db:"request_id"`
	Index      int    `json:"idx" db:"idx"`
	Status     int    `json:"status" db:"status"`
	Url        string `json:"url" db:"url"`
	Method     string `json:"method" db:"method"`
	Param      string `json:"param" db:"param"`
	Result     string `json:"try_result" db:"try_result"`
	CreateTime int64  `json:"create_time" db:"create_time"`
	UpdateTime int64  `json:"update_time" db:"update_time"`
	Resp       *util.Response
}
