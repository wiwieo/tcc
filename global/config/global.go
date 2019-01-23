package config

import (
	"flag"
	"tcc_transaction/store/data"
)

var (
	// 数据库连接
	C data.DataClient

	// email
	EmailUsername      = flag.String("email-username", "", "")
	EmailTo            = flag.String("email-to", "", "email receiver, if have many please use ',' split.")
	EmailPassword      = flag.String("email-password", "", "")
	MaxExceptionalData = flag.Int("max-exceptional-data", 100, "send msg when exceptional data than the value")

	// log
	LogFilePath = flag.String("log-file-path", "", "log file path")
	LogLevel    = flag.String("log-level", "", "log level")

	// mysql
	MysqlUsername = flag.String("mysql-username", "root", "")
	MysqlPassword = flag.String("mysql-password", "tcc_123", "")
	MysqlHost     = flag.String("mysql-host", "127.0.0.1", "")
	MysqlPort     = flag.String("mysql-port", "3306", "")
	MysqlDatabase = flag.String("mysql-database", "tcc", "")

	// levelDB
	DBPath = flag.String("level-db-path", "./tcc", "")

	// a interval time that is to execute task
	TimerInterval = flag.Int("timer-interval", 60*30, "unit is second")

	ApiKeyPrefix = flag.String("api-key-prefix", "/tcc/api/", "")
)
