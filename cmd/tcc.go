package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tcc_transaction/constant"
	"tcc_transaction/global/various"
	"tcc_transaction/log"
	"tcc_transaction/model"
	"tcc_transaction/store/data"
	"tcc_transaction/util"
	"time"
)

type tcc interface {
	// r：原生Request请求
	// api：根据当前请求，从配置文件中获取的Try的URL信息
	// 返回值：1、尝试过程中，成功的步骤
	// 2、错误信息
	Try(r *http.Request, api *model.RuntimeApi) ([]*data.SuccessStep, error)

	// r：原生Request请求
	// api：根据当前请求，从配置文件中获取的Confirm的URL信息
	// 返回值：1、错误信息
	Confirm(r *http.Request, api *model.RuntimeApi) error

	// r：原生Request请求
	// api：根据当前请求，从配置文件中获取的Cancel的URL信息
	// nodes：Try时可能成功的步骤，即需要回滚的步骤（根据Try返回值封装生成）
	// 返回值：1、执行取消时，失败步骤的ID编号集合
	// 2、错误信息
	Cancel(r *http.Request, api *model.RuntimeApi, nodes []*model.RuntimeTCC) ([]int64, error)
}

// 默认的处理逻辑
// 如果有和业务耦合无法剥离的情况，需要自定义处理
// 只要实现接口tcc的接口即可
type DefaultTcc struct {
}

func NewDefaultTcc() tcc {
	return &DefaultTcc{}
}

func (d *DefaultTcc) Try(r *http.Request, api *model.RuntimeApi) ([]*data.SuccessStep, error) {
	var success []*data.SuccessStep
	for _, node := range api.Nodes {
		var rst *util.Response
		tryURL := util.URLRewrite(api.UrlPattern, r.RequestURI[len(serverName)+1:], node.Try.Url)

		// try
		dt, err := util.HttpForward(tryURL, node.Try.Method, []byte(api.RequestInfo.Param), r.Header, time.Duration(node.Try.Timeout))
		// 不管成功与否（主要为了防止：当服务方接收并处理成功，但返回时失败），将结果保存起来，以备使用
		// 如果插入失败，则直接返回，并在后续回滚之前的步骤
		ss := &data.SuccessStep{
			RequestId: api.RequestInfo.Id,
			Index:     node.Index,
			Url:       tryURL,
			Method:    node.Try.Method,
			Param:     string(api.RequestInfo.Param),
			Result:    string(dt),
			Status:    constant.RequestTypeTry,
		}
		success = append(success, ss)

		if err != nil {
			log.Errorf("access try method failed, error info: %s", err)
			return success, err
		}

		err = json.Unmarshal(dt, &rst)
		ss.Resp = rst

		if err != nil {
			return success, err
		}

		if rst.Code != constant.Success {
			return success, fmt.Errorf(rst.Msg)
		}
	}
	return success, nil
}

func (d *DefaultTcc) Confirm(r *http.Request, api *model.RuntimeApi) error {
	for _, node := range api.Nodes {
		var rst *util.Response
		URL := util.URLRewrite(api.UrlPattern, r.RequestURI[len(serverName)+1:], node.Confirm.Url)

		// confirm or cancel
		dt, err := util.HttpForward(URL, node.Confirm.Method, []byte(api.RequestInfo.Param), r.Header, time.Duration(node.Confirm.Timeout))
		if err != nil {
			log.Errorf("confirm failed, please check it. error info is: %+v", err)
			return err
		}
		log.Infof("[%s] confirm response back content is: %+v", URL, string(dt))

		err = json.Unmarshal(dt, &rst)
		if err != nil {
			log.Errorf("confirm failed, please check it. error info is: %+v", err)
			return err
		}

		if rst.Code != constant.Success {
			err = fmt.Errorf(rst.Msg)
			log.Errorf("confirm failed, please check it. error info is: %+v", err)
			return err
		}

		// 处理成功后，修改状态
		various.C.Confirm(api.RequestInfo.Id)
	}

	// 全部提交成功，则修改状态为提交成功，避免重复调用
	various.C.UpdateRequestInfoStatus(constant.RequestInfoStatus1, api.RequestInfo.Id)

	return nil
}

func (d *DefaultTcc) Cancel(r *http.Request, api *model.RuntimeApi, nodes []*model.RuntimeTCC) ([]int64, error) {
	var ids []int64
	for _, node := range nodes {
		var rst *util.Response
		URL := util.URLRewrite(api.UrlPattern, r.RequestURI[len(serverName)+1:], node.Cancel.Url)

		// confirm or cancel
		dt, err := util.HttpForward(URL, node.Cancel.Method, []byte(api.RequestInfo.Param), r.Header, time.Duration(node.Cancel.Timeout))
		if err != nil {
			log.Errorf("cancel failed, please check it. error info is: %+v", err)
			return nil, err
		}
		log.Infof("[%s] cancel response back content is: %+v", URL, string(dt))

		err = json.Unmarshal(dt, &rst)
		if err != nil {
			return nil, err
		}

		if rst.Code != constant.Success {
			err = fmt.Errorf(rst.Msg)
			log.Errorf("cancel failed, please check it. error info is: %+v", err)
			return nil, err
		}

		// 如果当前数据有异常，则跳过此数据（交由后继异常流程处理）
		if node.SuccessStep.Id == 0 {
			continue
		}
		// 用于处理成功后，修改状态使用
		ids = append(ids, node.SuccessStep.Id)
	}
	return ids, nil
}
