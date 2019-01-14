package main

import (
	"fmt"
	"net/http"
	"tcc_transaction/constant"
	"tcc_transaction/global/various"
	"tcc_transaction/log"
	"tcc_transaction/model"
	"tcc_transaction/store/data"
	"tcc_transaction/util"
)

type proxy struct {
	t tcc
}

func (p *proxy) process(writer http.ResponseWriter, request *http.Request) {
	var response = &util.Response{}
	params := util.GetParams(request)
	log.Infof("welcome to tcc. url is %s, and param is %s", request.RequestURI, string(params))

	// 将请求信息持久化
	ri := &data.RequestInfo{
		Url:    request.RequestURI[len(serverName)+1:],
		Method: request.Method,
		Param:  string(params),
	}
	err := various.C.InsertRequestInfo(ri)
	if err != nil {
		response.Code = constant.InsertTccDataErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	}

	runtimeAPI, err := various.GetApiWithURL(request.RequestURI[len(serverName)+1:])
	if err != nil {
		response.Code = constant.NotFoundErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	}
	runtimeAPI.RequestInfo = ri

	// 转发--Try
	cancelSteps, err := p.try(request, runtimeAPI)

	if err != nil { // 回滚
		if len(cancelSteps) > 0 {
			go p.cancel(request, runtimeAPI, cancelSteps)
		}
		log.Errorf("try failed, error info is: %s", err)
		response.Code = constant.InsertTccDataErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	} else { // 提交
		go p.confirm(request, runtimeAPI)
	}
	response.Code = constant.Success
	util.ResponseWithJson(writer, response)
	return
}

func (p *proxy) try(r *http.Request, api *model.RuntimeApi) ([]*model.RuntimeTCC, error) {
	var nextCancelStep []*model.RuntimeTCC

	tryNodes := api.Nodes
	if len(tryNodes) == 0 {
		return nextCancelStep, fmt.Errorf("no method need to execute")
	}

	success, err := p.t.Try(r, api)
	if len(success) == 0 {
		return nil, fmt.Errorf("no success method")
	}

	err2 := various.C.BatchInsertSuccessStep(success)
	if err != nil || err2 != nil {
		for _, node := range api.Nodes {
			for _, s := range success {
				if node.Index == s.Index {
					node.SuccessStep = s
					nextCancelStep = append(nextCancelStep, node)
				}
			}
		}
		return nextCancelStep, err
	}
	if err2 != nil {
		return nextCancelStep, err2
	}
	return nextCancelStep, nil
}

func (p *proxy) confirm(r *http.Request, api *model.RuntimeApi) error {
	err := p.t.Confirm(r, api)
	if err != nil {
		various.C.UpdateRequestInfoStatus(constant.RequestInfoStatus2, api.RequestInfo.Id)
		return err
	}
	// 处理成功后，修改状态
	various.C.Confirm(api.RequestInfo.Id)
	// 全部提交成功，则修改状态为提交成功，避免重复调用
	various.C.UpdateRequestInfoStatus(constant.RequestInfoStatus1, api.RequestInfo.Id)
	return nil
}

func (p *proxy) cancel(r *http.Request, api *model.RuntimeApi, nodes []*model.RuntimeTCC) error {
	ids, err := p.t.Cancel(r, api, nodes)
	if err != nil {
		various.C.UpdateRequestInfoStatus(constant.RequestInfoStatus4, api.RequestInfo.Id)
		return err
	}
	for _, id := range ids {
		various.C.UpdateSuccessStepStatus(api.RequestInfo.Id, id, constant.RequestTypeCancel)
	}
	various.C.UpdateRequestInfoStatus(constant.RequestInfoStatus3, api.RequestInfo.Id)
	return nil
}
