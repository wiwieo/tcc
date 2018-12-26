package writer

import (
	"fmt"
	"tcc_transaction/log/writer/common"
	"tcc_transaction/log/writer/mmap"
	"tcc_transaction/log/writer/stdout"
)

// 如果有需要实现新的日志写入方式，则直接实现这个接口，并在创建实例时，返回对应的实例即可
type Writer interface {
	Write(content []byte) error
	Close() error
}

// 创建日志写入实例
func NewWriter(filePath string, size int) Writer {
	// 有指定日志文件，则使用mmap
	if len(filePath) > 1 {
		common.Mkdir(filePath)
		m, errM := mmap.NewMmap(filePath, size)
		// mmap映射失败，则使用终端
		if errM != nil {
			s, err := stdout.New()
			if err != nil {
				panic(err)
			}
			s.Write([]byte(fmt.Sprintf("mmap映射失败，使用终端日志。%s%s", errM, fmt.Sprintln())))
			return s
		}
		return m
	} else { // 路径不存在，则使用终端输出日志
		s, err := stdout.New()
		if err != nil {
			panic(err)
		}
		s.Write([]byte(fmt.Sprintf("mmap映射文件路径不存在，使用终端日志。%s", fmt.Sprintln())))
		return s
	}
}
