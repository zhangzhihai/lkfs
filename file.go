package lkfs

import (
	"fmt"
	//"io"
	"os"
)

const BlockSize = 1024 * 1024 * 64 //文件大小

type filesystem struct {
	path string
}

//找到当前写入的文件索引
//检查文件大小
func Newfs(path string) (fs *filesystem, err error) {

	fs = &filesystem{
		path: path,
	}
	return fs, nil
}

//文件的读取
//读写时独立使用变量避免加锁，使用独立的只读写只写句柄
//注意传指针
func (f *filesystem) Read(block int32, start int32, size int32) ([]byte, error) {

	fileName := fmt.Sprintf("%s/%d.dat", f.path, block)
	fi, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0)
	//定位指针
	fi.Seek(int64(start), 0)
	buf := make([]byte, size)
	fi.Read(buf)
	return buf, err
}

//写入前先检查文件大小
//文件大小写入全局中,使用写入锁
//注意打开关闭文件之文件大小
func (f *filesystem) Write(block int32, start int32, szie int32, body []byte) (bool, error) {

	fileName := fmt.Sprintf("%s/%d.dat", f.path, block)
	fi, err := os.OpenFile(fileName, os.O_RDWR, 0)
	//定位指针
	fi.Seek(int64(start), 0)

	_, err = fi.Write(body)
	if err != nil {
		return false, err
	}
	return true, nil
}
