package global

import (
	"flag"
	"fmt"
	"regexp"
	"tcc_transaction/log"
	"tcc_transaction/model"
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

	MysqlUsername = flag.String("mysql-username", "tcc", "")
	MysqlPassword = flag.String("mysql-password", "tcc_123", "")
	MysqlHost = flag.String("mysql-host", "localhost", "")
	MysqlPort = flag.String("mysql-port", "3306", "")
	MysqlDatabase = flag.String("mysql-database", "tcc", "")

	TimerInterval = flag.Int("timer-interval", 60 * 30, "unit is second")

	Apis = []*model.Api{
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
	}
)

func init() {
	flag.Parse()
	C = mysql.NewMysqlClient(*MysqlUsername, *MysqlPassword, *MysqlHost, *MysqlPort, *MysqlDatabase)

	log.InitLogrus(*LogFilePath, *LogLevel)
}

func GetApiWithURL(url string) (*model.RuntimeApi, error) {
	for _, v := range Apis {
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
