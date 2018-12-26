// +build !linux,!darwin

package mmap

import (
	"fmt"
)

//func MmapRead(filePath string) (content []byte, err error) {
//	// open a file
//	fd, err := syscall.Open(filePath, syscall.GENERIC_ALL, 0)
//	defer syscall.Close(fd)
//	if err != nil {
//		return nil, err
//	}
//	// get file size
//	fsize, err := syscall.Seek(fd, 0, 2)
//	if err != nil {
//		return nil, err
//	}
//	content = make([]byte, fsize)
//
//	// create memory mapping to file
//	hd, err := syscall.CreateFileMapping(fd, nil, syscall.PAGE_READONLY, 0, 0, nil)
//	if err != nil {
//		return nil, err
//	}
//	defer syscall.CloseHandle(hd)
//
//	// read file content
//	addr, err := syscall.MapViewOfFile(hd, syscall.FILE_MAP_READ, 0, 0, uintptr(fsize))
//	if err != nil {
//		return nil, err
//	}
//	// 释放映射的内存
//	defer syscall.UnmapViewOfFile(addr)
//	x := (*[BUFFER_SIZE]byte)(unsafe.Pointer(addr))
//	// 因为地址在出此方法会被释放，此处需要将内容拷贝到一个新的内存地址中
//	copy(content, x[:fsize])
//	return
//}
//
//func MmapWrite(filePath string, content []byte) error {
//	fd, err := syscall.Open(filePath, syscall.O_RDWR, syscall.O_RDWR)
//	defer syscall.Close(fd)
//	if err != nil {
//		return err
//	}
//	// get file size
//	fsize, err := syscall.Seek(fd, 0, 2)
//	if err != nil {
//		return err
//	}
//	// create memory mapping to file
//	hd, err := syscall.CreateFileMapping(fd, nil, syscall.PAGE_READWRITE, 0, 0, nil)
//	if err != nil {
//		return err
//	}
//	defer syscall.CloseHandle(hd)
//
//	// read file content
//	addr, err := syscall.MapViewOfFile(hd, syscall.FILE_MAP_WRITE, 0, 0, uintptr(fsize))
//	if err != nil {
//		return err
//	}
//	// 释放映射的内存
//	defer syscall.UnmapViewOfFile(addr)
//	x := (*[BUFFER_SIZE]byte)(unsafe.Pointer(addr))
//	println(fsize)
//	println(fmt.Sprintf("%+v", string(x[:len(content)])))
//	for i, v := range content {
//		x[i] = v
//	}
//	return nil
//}
type mmap struct {
}

func NewMmap(filePath string, size int) (*mmap, error) {
	return nil, fmt.Errorf("not support windows system")
}

func (n *mmap) Write(content []byte) error {
	return fmt.Errorf("not support windows system")
}

func (n *mmap) Close() error {
	return fmt.Errorf("not support windows system")
}
