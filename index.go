package lkfs

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"ulucu.github.com/log.v1"
)

/**
 * 索引系统
 * 检查系统
 **/

// 开始位置记录信息为128字节不足就补空,maxid int 2^32 4字节,BlockId 4字节,
//
type Index struct {
	MaxId     int32
	BlockId   int32
	Fi        *os.File
	ReadFile  *os.File
	WriteFile *os.File
}

const MaxIdx = 1024 * 1024 * 64 //字节数
const Ext = "idx"
const Hder = 128

//启动初始化对像
func NewIndex(path string) (idx *Index, err error) {
	//检查目录是否存在
	pt, err := Init(path)
	if err != nil {
		return nil, err
	}

	fmt.Println("pat", pt)
	fi, err := os.OpenFile(pt, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	//fmt.Println(pt)
	//初始化变量
	buf := make([]byte, 128)
	r := bufio.NewReader(fi)

	n, err := r.Read(buf)
	if n == 0 {
		return nil, errors.New("reader errors len 0")
	}

	var MaxId int32
	var BlockId int32

	binary.Read(bytes.NewBuffer(buf[0:3]), binary.BigEndian, &MaxId)
	binary.Read(bytes.NewBuffer(buf[4:7]), binary.BigEndian, &BlockId)

	idx = &Index{
		MaxId:   MaxId,
		BlockId: BlockId,
		Fi:      fi,
	}

	//fi.Close()

	return
}

func Init(path string) (pt string, err error) {
	b, err := PathExists(path)
	if err != nil {
		log.Fatalf("index init is path err", err)
		return "", err
	}

	//建立文件夹
	if b == false {
		m, err := PathMkdirAll(path)
		if err != nil {
			log.Fatalf("make path fail is path err", err)
			return "", err
		}
		if m == false {
			return "", errors.New("make dir err")
		}
	}

	//检查文件
	dir_list, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}

	for _, v := range dir_list {

		tn := v.Name()

		if strings.Contains(tn, Ext) == true {
			return fmt.Sprintf("%s%s%s", path, "/", v.Name()), nil
		}
	}

	//建立索引文件
	name := fmt.Sprintf("%s%s%d.%s", path, "/", 0, Ext)

	//如何把int 转成byte 存到文件中
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(0))
	//创建能接收buf写入的对像
	byt := bytes.NewBuffer(make([]byte, Hder-8))

	//定义个长度为128的buf

	byt.Write(buf)
	byt.Write(buf)

	//var d1 = []byte(wireteString)
	err = ioutil.WriteFile(name, byt.Bytes(), 0666)
	return name, err
}

//生成索引id
func (i *Index) Uuid() (id int32, err error) {
	i.MaxId = i.MaxId + 1
	return i.MaxId, nil
}

//生成12字节的0,1234,5678{blockid,开始位置，长度}
//写入索引
func (i *Index) Write(id int32, blockid int32, start int32, strlen int32) error {

	w := id*12 + Hder

	fmt.Println("w", w)

	i.Fi.Seek(int64(w), 0)

	var bufblockid, bufstart, bufstrlen []byte //= make([]byte, 4)

	bufblockid, bufstart, bufstrlen = make([]byte, 4), make([]byte, 4), make([]byte, 4)

	binary.BigEndian.PutUint32(bufblockid, uint32(blockid))

	binary.BigEndian.PutUint32(bufstart, uint32(start))

	binary.BigEndian.PutUint32(bufstrlen, uint32(strlen))

	byt := bytes.NewBuffer([]byte{})

	byt.Write(bufblockid)
	byt.Write(bufstart)
	byt.Write(bufstrlen)
	//写入文件
	fmt.Println("byt:", byt.Bytes())
	_, err := i.Fi.Write(byt.Bytes())

	cur_offset, _ := i.Fi.Seek(0, os.SEEK_CUR)
	fmt.Println("cur_offset", cur_offset)

	if err != nil {
		return err
	}

	//i.Fi.Close()
	return nil
}

func (i *Index) Reader(id int32) (blockid int32, start int32, strlen int32, err error) {
	w := id*12 + Hder

	fmt.Println("w", w)

	i.Fi.Seek(int64(w), 0)

	buf := make([]byte, 12)
	r := bufio.NewReader(i.Fi)

	n, err := r.Read(buf)
	if n == 0 {
		return 0, 0, 0, errors.New("reader errors len 0")
	}
	fmt.Println(buf)

	binary.Read(bytes.NewBuffer(buf[0:4]), binary.BigEndian, &blockid)
	binary.Read(bytes.NewBuffer(buf[4:8]), binary.BigEndian, &start)
	binary.Read(bytes.NewBuffer(buf[8:12]), binary.BigEndian, &strlen)
	return blockid, start, strlen, nil
}
