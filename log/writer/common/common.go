package common

import (
	"os"
	"strings"
	"time"
)

const CACHE_COUNT = 100

func GetTimeer(now time.Time) int64 {
	dest := time.Date(now.Year(), now.Month(), now.Day()+1, 1, 5, 0, 0, time.Local)
	return dest.UnixNano() - now.UnixNano()
}

func Mkdir(path string) {
	if flg, _ := pathExists(path); !flg {
		os.MkdirAll(path[:strings.LastIndex(path, string(os.PathSeparator))], os.ModePerm)
	}
}

// 判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
