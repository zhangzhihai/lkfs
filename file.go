package lkfs

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"ulucu.github.com/log.v1"
)

const BlockSize = 1024 * 1024 * 64 //文件大小
const Dxt = "dat"

type filesystem struct {
	path string
	idx  int32
	w    *os.File
	//r     *os.File
	mutex sync.Mutex
}

//找到当前写入的文件索引
//检查文件大小
//path 是文件路径
//idx 是最新的文件
func Newfs(path string, idx int32) (fs *filesystem, err error) {

	//检查文件夹是否存在
	log.Debugf("madir path", path)
	b, err := PathExists(path)
	if err != nil {
		return nil, err
	}

	if b == false {

		fmt.Printf("\n", path)

		m, err := PathMkdirAll(path)

		if err != nil {
			log.Fatalf("make path fail is path err", err)
			return nil, err
		}

		if m == false {
			return nil, errors.New("make dir err")
		}
	}

	fileName := fmt.Sprintf("%s/%d.%s", path, idx, Dxt)

	fi, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0664)

	//fr, err := os.OpenFile(fileName, os.O_RDONLY, 0)

	if err != nil {
		return nil, err
	}

	fs = &filesystem{
		path: path,
		idx:  idx,
		w:    fi,
		//r:    fr,
	}

	return fs, nil
}

//文件的读取
//读写时独立使用变量避免加锁，使用独立的只读写只写句柄
//注意传指针
func (f *filesystem) Read(block int32, start int32, size int32) ([]byte, error) {

	fileName := fmt.Sprintf("%s/%d.%s", f.path, block, Dxt)

	fi, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	//定位指针

	fi.Seek(int64(start), 0)

	buf := make([]byte, size)
	fi.Read(buf)
	fi.Close()

	return buf, err
}

//写入前先检查文件大小
//文件大小写入全局中,使用写入锁
//注意打开关闭文件之文件大小
func (f *filesystem) Write(szie int32, body []byte) (start int32, bid int32, err error) {

	var end int32 = 0

	fileName := fmt.Sprintf("%s/%d.%s", f.path, f.idx, Dxt)

	finfo, err := os.Stat(fileName)

	if err != nil {
		return end, f.idx, err
	}
	var n int64 = 0

	if finfo.Size() > BlockSize {

		f.w.Close()

		f.mutex.Lock()

		f.idx += 1
		fileName = fmt.Sprintf("%s/%d.%s", f.path, f.idx, Dxt)

		f.w, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0664)
		if err != nil {
			return end, f.idx, err
		}
		//f.w.Seek(int64(0), 0)

	} else {

		f.mutex.Lock()
		// 查找文件末尾的偏移量
		n, err = f.w.Seek(0, os.SEEK_END)
		if err != nil {
			return end, f.idx, nil
		}
		//f.w.Seek(int64(n), 0)
	}

	_, err = f.w.WriteAt(body, n)
	f.mutex.Unlock()

	if err != nil {
		return end, f.idx, err
	}
	return int32(n), f.idx, nil
}
