package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"tcc_transaction/log"
	"time"
)

func HttpForward(url, method string, param []byte, head map[string][]string, timeout time.Duration) ([]byte, error) {
	log.Infof("url: %s, method: %s, param: %s", url, method, string(param))
	switch method {
	case http.MethodGet:
		rst, err := httpGet(url, head, timeout)
		log.Infof("url: %s, method: %s, param: %s, response: %s, error: %s", url, method, string(param), string(rst), err)
		return rst, err
	case http.MethodPost:
		rst, err := httpPost(url, head, timeout, param)
		log.Infof("url: %s, method: %s, param: %s, response: %s, error: %s", url, method, string(param), string(rst), err)
		return rst, err
	default:
		return nil, fmt.Errorf("not support this method, please implement it")
	}
}

func httpPost(url string, head map[string][]string, timeout time.Duration, param []byte) ([]byte, error) {
	if timeout == 0 {
		timeout = time.Hour
	}
	r, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(param))
	if err != nil {
		return nil, err
	}

	setHead(r, head)

	rsp, err := (&http.Client{
		Timeout: timeout,
	}).Do(r)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(rsp.Status)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func httpGet(url string, head map[string][]string, timeout time.Duration) ([]byte, error) {
	if timeout == 0 {
		timeout = time.Hour
	}
	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	setHead(r, head)

	rsp, err := (&http.Client{
		Timeout: timeout,
	}).Do(r)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(rsp.Status)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func setHead(r *http.Request, head map[string][]string) {
	if len(head) == 0 {
		return
	}
	for k, v := range head {
		r.Header.Add(k, v[0])
	}
}

func GetParams(r *http.Request) []byte {
	switch r.Method {
	case http.MethodGet:
		p := r.URL.Query()
		var params = make(map[string]interface{})
		for k, v := range p {
			params[k] = v[0]
		}
		rst, err := json.Marshal(params)
		if err != nil {
			panic(err)
		}
		return rst
	case http.MethodPost, http.MethodPut:
		p, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		return p
	default:
		return nil
	}
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseWithJson(write http.ResponseWriter, r *Response) {
	write.Header().Set("ContentType", "Application/json")
	content, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	write.Write(content)
}

// url rewrite
// "/api/(.*)/actions/(.*)" to "/api/v1/$1/actions/$2"
func URLRewrite(partern, o, d string) string {
	reg, err := regexp.Compile(partern)
	if err != nil {
		return d
	}
	return reg.ReplaceAllString(o, d)
}
