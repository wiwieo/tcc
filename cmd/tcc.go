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
	"tcc_transaction/util"
	"time"
)


func tcc(writer http.ResponseWriter, request *http.Request) {
	var response = &util.Response{}
	params := util.GetParams(request)
	log.Infof("welcome to tcc. url is %s, and param is %s", request.RequestURI, string(params))

	// 将请求信息持久化
	ri := &data.RequestInfo{
		Url:    request.RequestURI[len(serverName)+1:],
		Method: request.Method,
		Param:  string(params),
	}
	err := global.C.InsertRequestInfo(ri)
	if err != nil {
		response.Code = constant.InsertTccDataErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	}

	runtimeAPI, err := global.GetApiWithURL(request.RequestURI[len(serverName)+1:])
	if err != nil {
		response.Code = constant.NotFoundErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	}
	runtimeAPI.RequestInfo = ri

	// 转发--Try
	cancelSteps, err := try(request, runtimeAPI)

	if err != nil { // 回滚
		if len(cancelSteps) > 0 {
			go cancel(request, runtimeAPI, cancelSteps)
		}
		log.Errorf("try failed, error info is: %s", err)
		response.Code = constant.InsertTccDataErrCode
		response.Msg = err.Error()
		util.ResponseWithJson(writer, response)
		return
	} else { // 提交
		go confirm(request, runtimeAPI)
	}
	response.Code = constant.Success
	util.ResponseWithJson(writer, response)
	return
}

func try(r *http.Request, api *model.RuntimeApi) ([]*model.RuntimeTCC, error) {
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
		err = global.C.InsertSuccessStep(ss)
		tryNodes[idx].SuccessStep = ss
		nextCancelStep = append(nextCancelStep, tryNodes[idx])
		if err != nil {
			log.Errorf("insert into success_step failed, need to special process, error info: %s", err)
			return nextCancelStep, err
		}
	}
	return nil, nil
}

func confirm(r *http.Request, api *model.RuntimeApi) error {
	var err error
	for _, node := range api.Nodes {
		var rst *util.Response
		URL := util.URLRewrite(api.UrlPattern, r.RequestURI[len(serverName)+1:], node.Confirm.Url)

		// confirm or cancel
		dt, err := util.HttpForward(URL, node.Confirm.Method, []byte(api.RequestInfo.Param), r.Header, time.Duration(node.Confirm.Timeout))
		if err != nil {
			goto ERROR
		}
		log.Infof("[%s] confirm response back content is: %+v", URL, string(dt))

		err = json.Unmarshal(dt, &rst)
		if err != nil {
			goto ERROR
		}

		if rst.Code != constant.Success {
			err = fmt.Errorf(rst.Msg)
			goto ERROR
		}

		// 处理成功后，修改状态
		global.C.Confirm(api.RequestInfo.Id)
	}

	// 全部提交成功，则修改状态为提交成功，避免重复调用
	global.C.UpdateRequestInfoStatus(constant.RequestInfoStatus_1, api.RequestInfo.Id)

	return nil
ERROR:
	global.C.UpdateRequestInfoStatus(constant.RequestInfoStatus_2, api.RequestInfo.Id)
	log.Errorf("confirm failed, please check it. error info is: %+v", err)
	return err
}

func cancel(r *http.Request, api *model.RuntimeApi, nodes []*model.RuntimeTCC) error {
	var err error
	for _, node := range nodes {
		var rst *util.Response
		URL := util.URLRewrite(api.UrlPattern, r.RequestURI[len(serverName)+1:], node.Cancel.Url)

		// confirm or cancel
		dt, err := util.HttpForward(URL, node.Cancel.Method, []byte(api.RequestInfo.Param), r.Header, time.Duration(node.Cancel.Timeout))
		if err != nil {
			goto ERROR
		}
		log.Infof("[%s] cancel response back content is: %+v", URL, string(dt))

		err = json.Unmarshal(dt, &rst)
		if err != nil {
			goto ERROR
		}

		if rst.Code != constant.Success {
			err = fmt.Errorf(rst.Msg)
			goto ERROR
		}

		// 如果当前数据有异常，则跳过此数据（交由后继异常流程处理）
		if node.SuccessStep.Id == 0 {
			continue
		}
		// 处理成功后，修改状态
		global.C.UpdateSuccessStepStatus(node.SuccessStep.Id, constant.RequestTypeCancel)
	}
	global.C.UpdateRequestInfoStatus(constant.RequestInfoStatus_3, api.RequestInfo.Id)
	return nil
ERROR:
	log.Errorf("cancel failed, please check it. error info is: %+v", err)
	global.C.UpdateRequestInfoStatus(constant.RequestInfoStatus_4, api.RequestInfo.Id)
	return err
}
