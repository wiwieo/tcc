package model

import "tcc_transaction/store/data"

type Api struct {
	UrlPattern string `json:"url_pattern"`
	Nodes      []*TCC `json:"nodes"`
}

type TCC struct {
	Index   int   `json:"index"`
	Try     *Node `json:"try"`
	Confirm *Node `json:"confirm"`
	Cancel  *Node `json:"cancel"`
}

type Node struct {
	Url     string `json:"url"`
	Method  string `json:"method"`
	Timeout int    `json:"timeout"`
	Param   string `json:"param"` // 暂时先不使用
}

type RuntimeApi struct {
	UrlPattern  string
	RequestInfo *data.RequestInfo
	Nodes       []*RuntimeTCC
}

type RuntimeTCC struct {
	Index       int
	Try         *RuntimeNode
	Confirm     *RuntimeNode
	Cancel      *RuntimeNode
	SuccessStep *data.SuccessStep
}

type RuntimeNode struct {
	Url     string
	Method  string
	Timeout int
}

func ConverToRuntime(nodes []*TCC) []*RuntimeTCC {
	var rns = make([]*RuntimeTCC, 0, len(nodes))
	for _, n := range nodes {
		rns = append(rns, &RuntimeTCC{
			Index: n.Index,
			Try: &RuntimeNode{
				Url:     n.Try.Url,
				Method:  n.Try.Method,
				Timeout: n.Try.Timeout,
			},
			Confirm: &RuntimeNode{
				Url:     n.Confirm.Url,
				Method:  n.Confirm.Method,
				Timeout: n.Confirm.Timeout,
			},
			Cancel: &RuntimeNode{
				Url:     n.Cancel.Url,
				Method:  n.Cancel.Method,
				Timeout: n.Cancel.Timeout,
			},
		})
	}
	return rns
}
