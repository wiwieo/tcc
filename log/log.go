// Copyright (c) 2018, dmc (814172254@qq.com),
//
// Authors: dmc,
//
// Distribution:.
package log

// InitLogrus 初始化日志配置
func InitLogrus(logPath, level string) error {
	SetPath(logPath)
	SetLevel(level)
	return nil
}
