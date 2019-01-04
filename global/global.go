package global

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"tcc_transaction/log"
	"tcc_transaction/model"
	"tcc_transaction/store/config/etcd"
	"tcc_transaction/store/data"
	"tcc_transaction/store/data/mysql"
	"time"
)

var (
	// 数据库连接
	C data.DataClient

	EmailUsername      = flag.String("email-username", "", "")
	EmailTo            = flag.String("email-to", "", "email receiver, if have many please use ',' split.")
	EmailPassword      = flag.String("email-password", "", "")
	MaxExceptionalData = flag.Int("max-exceptional-data", 100, "send msg when exceptional data than the value")

	LogFilePath = flag.String("log-file-path", "", "log file path")
	LogLevel    = flag.String("log-level", "", "log level")

	MysqlUsername = flag.String("mysql-username", "root", "")
	MysqlPassword = flag.String("mysql-password", "tcc_123", "")
	MysqlHost     = flag.String("mysql-host", "127.0.0.1", "")
	MysqlPort     = flag.String("mysql-port", "3306", "")
	MysqlDatabase = flag.String("mysql-database", "tcc", "")

	TimerInterval = flag.Int("timer-interval", 60*30, "unit is second")

	ApiKeyPrefix = flag.String("api-key-prefix", "/tcc/api/", "")
)

var (
	apis     []*model.Api
	etcdC, _ = etcd3.NewEtcd3Client([]string{"localhost:2379"}, int(time.Minute), "", "", nil)
)

func InitAll() {
	flag.Parse()
	C = mysql.NewMysqlClient(*MysqlUsername, *MysqlPassword, *MysqlHost, *MysqlPort, *MysqlDatabase)

	log.InitLogrus(*LogFilePath, *LogLevel)

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

	data, err := etcdC.List(context.Background(), *ApiKeyPrefix)
	if err != nil {
		panic(err)
	}
	apis = parseToApi(data)
	println(len(apis))
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
			if v.UrlPattern == string(k)[len(*ApiKeyPrefix):] {
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

	etcdC.WatchTree(context.Background(), *ApiKeyPrefix, callback)
}
