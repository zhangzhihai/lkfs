package lkfs

import (
	"fmt"
	"os"
	//"path/filepath"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathMkdirAll(path string) (bool, error) {
	dir, err := os.Getwd()
	if err != nil {
		return false, err
	}

	err = os.MkdirAll(fmt.Sprintf("%s%s", dir, path), os.ModePerm)
	if err != nil {
		return false, err
	}
	return true, nil
}

//检查文件大小
//写入大文件检查
