package main

import (
	"context"
	"encoding/json"
	"tcc_transaction/global/config"
	"tcc_transaction/model"
	"tcc_transaction/store/config/etcd"
	"time"
)

var apis = []*model.Api{
	{
		UrlPattern: "^accounts/order/(.)*",
		Nodes: []*model.TCC{
			{
				Index: 0,
				Try: &model.Node{
					Url:     "http://localhost:8083/accounts/order/try/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
				Confirm: &model.Node{
					Url:     "http://localhost:8083/accounts/order/confirm/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
				Cancel: &model.Node{
					Url:     "http://localhost:8083/accounts/order/cancel/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
			}, {
				Index: 1,
				Try: &model.Node{
					Url:     "http://localhost:8084/orders/order/try/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
				Confirm: &model.Node{
					Url:     "http://localhost:8084/orders/order/confirm/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
				Cancel: &model.Node{
					Url:     "http://localhost:8084/orders/order/cancel/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
			},
		},
	},
	{
		UrlPattern: "^examples/(.)*",
		Nodes: []*model.TCC{
			{
				Index: 0,
				Try: &model.Node{
					Url:     "http://localhost:8083/accounts/order/try/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
				Confirm: &model.Node{
					Url:     "http://localhost:8083/accounts/order/confirm/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
				Cancel: &model.Node{
					Url:     "http://localhost:8083/accounts/order/cancel/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
			}, {
				Index: 1,
				Try: &model.Node{
					Url:     "http://localhost:8084/orders/order/try/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
				Confirm: &model.Node{
					Url:     "http://localhost:8084/orders/order/confirm/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
				Cancel: &model.Node{
					Url:     "http://localhost:8084/orders/order/cancel/$1",
					Method:  "POST",
					Timeout: 5 * int(time.Second),
				},
			},
		},
	},
}

func main() {
	put()
}

func put() {
	var st, err = etcd3.NewEtcd3Client([]string{"localhost:2379"}, int(time.Minute), "", "", nil)
	if err != nil {
		panic(err)
	}
	for _, v := range apis {
		// 简单点，使用json，proto buff太麻烦了
		data, err := json.Marshal(v)
		if err != nil {
			panic(err)
		}
		err = st.Put(context.Background(), *config.ApiKeyPrefix+v.UrlPattern, data, 0)
		if err != nil {
			panic(err)
		}
	}
}
