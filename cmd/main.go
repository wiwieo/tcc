package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tcc_transaction/constant"
	"tcc_transaction/global"
	"tcc_transaction/log"
	"tcc_transaction/model"
	"tcc_transaction/store/data"
	"tcc_transaction/store/data/mysql"
	"tcc_transaction/task"
	"tcc_transaction/util"
	"time"
)

var serverName = "/tcc"

type GlobalInfo struct {
	c *mysql.MysqlClient
}

func init() {
	log.InitLogrus("", "")
	go task.TimerToExcuteTask()
}

func InitInfo() *GlobalInfo {
	c := mysql.NewMysqlClient("tcc", "tcc_123", "localhost", "3306", "tcc")
	return &GlobalInfo{
		c: c,
	}
}

func main() {
	gi := InitInfo()
	http.Handle("/", http.FileServer(http.Dir("file")))
	http.HandleFunc(fmt.Sprintf("%s/", serverName), gi.tcc)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func (gi *GlobalInfo) tcc(writer http.ResponseWriter, request *http.Request) {
	var response = &util.Response{}
	params := util.GetParams(request)
	log.Infof("welcome to tcc. url is %s, and param is %s", request.RequestURI, string(params))

	// 将请求信息持久化
	ri := &data.RequestInfo{
		Url:    request.RequestURI[len(serverName)+1:],
		Method: request.Method,
		Param:  string(params),
	}
	id, err := gi.c.InsertRequestInfo(ri)
	ri.Id = id
	if err != nil {
		response.Code = constant.InsertTccDataErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	}

	// 转发--Try
	runtimeAPI, err := global.GetApiWithURL(request.RequestURI[len(serverName)+1:])
	if err != nil {
		response.Code = constant.NotFoundErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	}
	runtimeAPI.RequestInfo = ri
	nextSteps, err := gi.try(request, runtimeAPI)

	if len(nextSteps) == 0 {
		response.Code = constant.NotFoundErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	}

	if err != nil { // 回滚
		go func() {
			err := gi.cancel(request, runtimeAPI.UrlPattern, []byte(runtimeAPI.RequestInfo.Param), nextSteps)
			if err != nil {
				log.Errorf("cancel failed, please check it. error info is: %+v", err)
				gi.c.UpdateRequestInfo(constant.RequestInfoStatus_4, id)
			} else {
				gi.c.UpdateRequestInfo(constant.RequestInfoStatus_3, id)
			}
		}()
		log.Errorf("try failed, error info is: %s", err)
		response.Code = constant.InsertTccDataErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	} else { // 提交
		go func() {
			err := gi.confirm(request, runtimeAPI)
			if err != nil {
				gi.c.UpdateRequestInfo(constant.RequestInfoStatus_2, id)
				log.Errorf("confirm failed, please check it. error info is: %+v", err)
			} else {
				gi.c.UpdateRequestInfo(constant.RequestInfoStatus_1, id)
			}
		}()
	}
	response.Code = constant.Success
	util.ResponseWithJson(writer, response)
	return
}

func (gi *GlobalInfo) try(r *http.Request, api *model.RuntimeApi) ([]*model.RuntimeTCC, error) {
	var nextCancelStep []*model.RuntimeTCC

	tryNodes := api.Nodes
	if len(tryNodes) == 0 {
		return nextCancelStep, fmt.Errorf("no method need to execute")
	}
	for idx, node := range tryNodes {
		var rst *util.Response
		tryURL := util.URLRewrite(api.UrlPattern, r.RequestURI[len(serverName)+1:], node.Try.Url)

		// try
		dt, err := util.HttpForward(tryURL, node.Try.Method, []byte(api.RequestInfo.Param), r.Header, time.Duration(node.Try.Timeout))
		if err != nil {
			return nextCancelStep, err
		}
		log.Infof("[%s] try response back content is: %+v", tryURL, string(dt))

		err = json.Unmarshal(dt, &rst)
		if err != nil {
			return nextCancelStep, err
		}

		if rst.Code != constant.Success {
			return nextCancelStep, fmt.Errorf(rst.Msg)
		}
		// 成功之后，将结果保存起来，以备使用
		// TODO：如果插入失败，则无法处理，需要人工干预
		ss := &data.SuccessStep{
			RequestId: api.RequestInfo.Id,
			Index:     node.Index,
			Url:       tryURL,
			Method:    node.Try.Method,
			Param:     string(api.RequestInfo.Param),
			Result:    string(dt),
			Status:    constant.RequestTypeTry,
		}
		sid, err := gi.c.InsertSuccessStep(ss)
		ss.Id = sid
		tryNodes[idx].SuccessStep = ss
		nextCancelStep = append(nextCancelStep, tryNodes[idx])
		if err != nil {
			log.Errorf("insert into success_step failed, need to special process, error info: %s", err)
			return nextCancelStep, err
		}
	}
	return nil, nil
}

func (gi *GlobalInfo) confirm(r *http.Request, api *model.RuntimeApi) error {
	for _, node := range api.Nodes {
		var rst *util.Response
		URL := util.URLRewrite(api.UrlPattern, r.RequestURI[len(serverName)+1:], node.Confirm.Url)

		// confirm or cancel
		dt, err := util.HttpForward(URL, node.Confirm.Method, []byte(api.RequestInfo.Param), r.Header, time.Duration(node.Confirm.Timeout))
		if err != nil {
			return err
		}
		log.Infof("[%s] confirm response back content is: %+v", URL, string(dt))

		err = json.Unmarshal(dt, &rst)
		if err != nil {
			return err
		}

		if rst.Code != constant.Success {
			return fmt.Errorf(rst.Msg)
		}

		// 处理成功后，修改状态
		gi.c.Confirm(api.RequestInfo.Id)
	}
	return nil
}

func (gi *GlobalInfo) cancel(r *http.Request, key string, param []byte, nodes []*model.RuntimeTCC) error {
	for _, node := range nodes {
		var rst *util.Response
		URL := util.URLRewrite(key, r.RequestURI[len(serverName)+1:], node.Cancel.Url)

		// confirm or cancel
		dt, err := util.HttpForward(URL, node.Cancel.Method, param, r.Header, time.Duration(node.Cancel.Timeout))
		if err != nil {
			return err
		}
		log.Infof("[%s] cancel response back content is: %+v", URL, string(dt))

		err = json.Unmarshal(dt, &rst)
		if err != nil {
			return err
		}

		if rst.Code != constant.Success {
			return fmt.Errorf(rst.Msg)
		}

		// 处理成功后，修改状态
		gi.c.UpdateSuccessStepStatus(node.SuccessStep.Id, constant.RequestTypeCancel)
	}
	return nil
}

