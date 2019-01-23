package various

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"tcc_transaction/global/config"
	"tcc_transaction/log"
	"tcc_transaction/model"
	"tcc_transaction/store/config/etcd"
	"tcc_transaction/store/data"
	"tcc_transaction/store/data/leveldb"
	"tcc_transaction/store/data/mysql"
	"time"
)

var (
	// 数据库连接
	C        data.DataClient
	apis     []*model.Api
	EtcdC, _ = etcd3.NewEtcd3Client([]string{"localhost:2379"}, int(time.Minute), "", "", nil)
)

func InitAll() {
	flag.Parse()
	var err error
	C, err = mysql.NewMysqlClient(*config.MysqlUsername, *config.MysqlPassword, *config.MysqlHost, *config.MysqlPort, *config.MysqlDatabase)
	if err != nil {
		C, err = leveldb.NewLevelDB(*config.DBPath)
		if err != nil {
			panic(err)
		}
	}
	log.InitLogrus(*config.LogFilePath, *config.LogLevel)

	LoadApiFromEtcd()
	WatchApi()
}

func GetApiWithURL(url string) (*model.RuntimeApi, error) {
	for _, v := range apis {
		reg, _ := regexp.Compile(v.UrlPattern)
		if reg.MatchString(url) {
			ra := &model.RuntimeApi{
				UrlPattern: v.UrlPattern,
				Nodes:      model.ConverToRuntime(v.Nodes),
			}
			return ra, nil
		}
	}
	return nil, fmt.Errorf("there is no api to run")
}

func LoadApiFromEtcd() {
	var parseToApi = func(param [][]byte) []*model.Api {
		var rsts = make([]*model.Api, 0, len(param))
		for _, bs := range param {
			var api *model.Api
			err := json.Unmarshal(bs, &api)
			if err != nil {
				panic(err)
			}
			rsts = append(rsts, api)
		}
		return rsts
	}

	data, err := EtcdC.List(context.Background(), *config.ApiKeyPrefix)
	if err != nil {
		panic(err)
	}
	apis = parseToApi(data)
}

func WatchApi() {
	del := func(idx int) {
		apis = append(apis[:idx], apis[idx+1:]...)
	}

	add := func(v []byte) error {
		var api *model.Api
		err := json.Unmarshal(v, &api)
		if err != nil {
			return err
		}
		apis = append(apis, api)
		return nil
	}

	modify := func(idx int, v []byte) error {
		del(idx)
		return add(v)
	}

	getIdx := func(k []byte) int {
		for idx, v := range apis {
			if v.UrlPattern == string(k)[len(*config.ApiKeyPrefix):] {
				return idx
			}
		}
		return 0
	}

	callback := func(k, v []byte, tpe string) {
		switch tpe {
		case etcd3.WatchEventTypeD:
			del(getIdx(k))
		case etcd3.WatchEventTypeC:
			add(v)
		case etcd3.WatchEventTypeM:
			modify(getIdx(k), v)
		}
	}

	EtcdC.WatchTree(context.Background(), *config.ApiKeyPrefix, callback)
}
