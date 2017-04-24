package lkfs

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"ulucu.github.com/log.v1"
)

/**
 * 索引系统
 * 检查系统
 **/

// 开始位置记录信息为128字节不足就补空,maxid int 2^32 4字节,BlockId 4字节,
//

const IndexkSize = 1024 * 1024 * 64 //文件大小

type FileIdx struct {
	MaxId   int32
	BlockId int32
	w       *os.File
	r       *os.File
	Mutex   sync.Mutex
}

//读
type ReaderIdx struct {
	r     *os.File
	stat  bool
	mutex sync.Mutex
}

var mr map[int]ReaderIdx

const Binext = "bin.idx" //索引文件名
const Hder = 128         //头长
const IdLen = 12

//启动初始化对像
func NewFildIdx(path string) (idx *FileIdx, err error) {
	mr = make(map[int]ReaderIdx)
	//检查目录是否存在
	pt, err := Init(path)

	if err != nil {
		return nil, err
	}

	fi, err := os.OpenFile(pt, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	ri, err := os.OpenFile(pt, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 128)
	r := bufio.NewReader(fi)

	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, errors.New("reader errors len 0")
	}

	var MaxId int32
	var BlockId int32

	binary.Read(bytes.NewBuffer(buf[0:3]), binary.BigEndian, &MaxId)
	binary.Read(bytes.NewBuffer(buf[4:7]), binary.BigEndian, &BlockId)

	idx = &FileIdx{
		MaxId:   MaxId,
		BlockId: BlockId,
		w:       fi,
		r:       ri,
	}

	//fmt.Println("idx", idx)
	return
}

//path
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

	//一个目录下只能有一个索引文件

	//建立索引文件
	name := fmt.Sprintf("%s/%s", path, Binext)

	//如何把int 转成byte 存到文件中
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(0))
	newbuf := bytes.NewBuffer(make([]byte, Hder-8))

	//定义个长度为128的buf
	newbuf.Write(buf)
	newbuf.Write(buf)
	//log.Println("init write len 128")
	err = ioutil.WriteFile(name, newbuf.Bytes(), 0666)
	return name, err
}

func (fi *FileIdx) Getr() (*os.File, int, error) {

	ml := 20
	i := len(mr)

	if len(mr) < ml {
		//log.Println("len mr i", i)
		for ; i < ml; i++ {
			//copy os.file

			mr[i] = ReaderIdx{
				r:    fi.r,
				stat: false,
			}
		}
	}

	//log.Println("len mr", len(mr))

	for k, v := range mr {
		if v.stat == false {
			v.stat = true
			mr[k] = v
			return mr[k].r, k, nil
		}
	}
	return fi.r, 0, nil
}

func (fi *FileIdx) Push(i int) error {
	v, ok := mr[i]
	if ok == false {
		delete(mr, i)
	} else {
		v.stat = false
		mr[i] = v
	}
	return nil
}

//生成索引id
func (fi *FileIdx) Uuid() (id int32, err error) {
	//fmt.Println("Maxid:", fi.MaxId)
	fi.MaxId = fi.MaxId + 1

	return fi.MaxId, nil
}

//生成12字节的0,1234,5678{blockid,开始位置，长度}
//写入索引
func (fi *FileIdx) Write(id int32, blockid int32, start int32, strlen int32) error {

	position := id*IdLen + Hder

	fi.Mutex.Lock()
	//移动写入指针到指定位置
	fi.w.Seek(int64(position), 0)

	var bufblockid, bufstart, bufstrlen []byte //= make([]byte, 4)

	bufblockid, bufstart, bufstrlen = make([]byte, 4), make([]byte, 4), make([]byte, 4)

	binary.BigEndian.PutUint32(bufblockid, uint32(blockid))
	binary.BigEndian.PutUint32(bufstart, uint32(start))
	binary.BigEndian.PutUint32(bufstrlen, uint32(strlen))

	bf := bytes.NewBuffer([]byte{})

	bf.Write(bufblockid)
	bf.Write(bufstart)
	bf.Write(bufstrlen)

	//写入文件
	//fmt.Println("byt:", byt.Bytes())
	_, err := fi.w.Write(bf.Bytes())

	//cur_offset, _ := fi.w.Seek(0, os.SEEK_CUR)
	//fmt.Println("cur_offset", cur_offset)
	fi.Mutex.Unlock()

	if err != nil {
		return err
	}

	return nil
}

func (fi *FileIdx) Reader(id int32) (blockid int32, start int32, strlen int32, err error) {

	position := id*IdLen + Hder
	f, i, err := fi.Getr()
	if err != nil {
		return 0, 0, 0, err
	}

	f.Seek(int64(position), 0)

	buf := make([]byte, IdLen)
	r := bufio.NewReader(f)

	n, err := r.Read(buf)
	if n == 0 {
		return 0, 0, 0, errors.New("reader errors len 0")
	}

	binary.Read(bytes.NewBuffer(buf[0:4]), binary.BigEndian, &blockid)
	binary.Read(bytes.NewBuffer(buf[4:8]), binary.BigEndian, &start)
	binary.Read(bytes.NewBuffer(buf[8:12]), binary.BigEndian, &strlen)
	err = fi.Push(i)
	if err != nil {
		return 0, 0, 0, errors.New("reader push errors len 0")
	}
	return blockid, start, strlen, nil
}
