package lkfs

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"ulucu.github.com/log.v1"
)

type mpso struct {
	path   string
	idx    int32
	dicmap map[string]string
	w      *os.File
	mutex  sync.Mutex
}

const dict = "dictionary.txt"   // 字典
const mpdict = "dictionmap.txt" //索引字典

const MpBlockSize = 1024 * 1024 * 64 //文件大小

//读取map,如果不存在读使用dict进行初始化
//map结构 {k:idx} map[string]int32,int32 (key,inxid,id)
//idx结构 int32,int32,int32,int32,int32(文章id,上一篇文章id位置,出现的次数，文章开始位置,文章结束位置)
func Newmap(path string) (f *mpso, err error) {

	b, err := PathExists(path)
	if err != nil {
		return nil, err
	}
	if b == false {
		log.Fatalf("dict data error ", err)
		return nil, err
	}

	dictfileName := fmt.Sprintf("%s/%s", path, dict)
	dictfb, err := PathExists(dictfileName)
	if dictfb == false {
		log.Fatalf("dict data exists err ", err)
		return nil, err
	}

	dictmpfileName := fmt.Sprintf("%s/%s", path, mpdict)
	dictmpfb, err := PathExists(dictmpfileName)

	var dicidxtmp map[string]string

	//定义map
	dicidxtmp = make(map[string]string)
	//定义idx
	var dIdx int32 = 0

	if dictmpfb == false {
		//生成字典map
		ra, err := ioutil.ReadFile(dictfileName)
		if err != nil {
			return nil, err
		}

		strarr := strings.Split(string(ra), "\n")
		//定义map

		for _, v := range strarr {
			astr := strings.Split(v, " ")
			k := astr[0]
			//dk := [1]int32{0}
			//fmt.Println(dk)
			dicidxtmp[k] = "0:0"
		}

		dicidxjson, err := json.Marshal(dicidxtmp)
		if err != nil {
			return nil, err
		}

		buf := make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(0))
		bbf := make([]byte, Hder-4)
		binary.BigEndian.PutUint32(bbf, uint32(0))

		newbuf := bytes.NewBuffer(make([]byte, 0, Hder+len(dicidxjson)))

		newbuf.Write(buf)
		newbuf.Write(bbf)

		newbuf.Write(dicidxjson)

		ioutil.WriteFile(dictmpfileName, newbuf.Bytes(), os.ModeAppend)

	} else {
		str, err := ioutil.ReadFile(dictmpfileName)
		//dIdx = str[0:4]
		binary.Read(bytes.NewBuffer(str[0:4]), binary.BigEndian, &dIdx)

		err = json.Unmarshal(str[128:], &dicidxtmp)
		if err != nil {
			return nil, err
		}

	}

	//dIdx
	fileName := fmt.Sprintf("%s/%d.bin", path, dIdx)

	//fmt.Println(fileName)

	fi, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0664)
	//fi, err := os.Open(fileName)
	/**
	bf := bytes.NewBuffer(make([]byte, 0, 1024))
	bf.Write([]byte("abcd1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"))
	x, err := fi.Write(bf.Bytes())
	**/

	f = &mpso{
		path:   path,
		idx:    dIdx,
		dicmap: dicidxtmp,
		w:      fi,
	}
	return f, nil
	//检查

}

/**
 * 查询 字典中的最新id与block文件,start 文件所以位置
 * @param  {[type]} m *mpso)        Get() (block, start int32, err error [description]
 * @return {[type]}   [description]
 **/

func (m *mpso) Get(block, start int32) (id, previous, count, lefstr, rightstr int32, err error) {

	file := fmt.Sprintf("%s/%d.bin", m.path, block)
	ri, err := os.Open(file)

	ri.Seek(int64(start), 0)

	buf := make([]byte, 20)

	ri.Read(buf)

	ri.Close()

	////fmt.Println(buf)

	binary.Read(bytes.NewBuffer(buf[0:4]), binary.BigEndian, &id)
	binary.Read(bytes.NewBuffer(buf[4:8]), binary.BigEndian, &previous)
	binary.Read(bytes.NewBuffer(buf[8:12]), binary.BigEndian, &count)
	binary.Read(bytes.NewBuffer(buf[12:16]), binary.BigEndian, &lefstr)
	binary.Read(bytes.NewBuffer(buf[16:20]), binary.BigEndian, &rightstr)
	//fmt.Printf("id %d,previous %d,count %d,lefstr %d,rightstr %d", id, previous, count, lefstr, rightstr)
	//log.Infof("id %d,previous %d,count %d,lefstr %d,rightstr %d", id, previous, count, lefstr, rightstr)
	return id, previous, count, lefstr, rightstr, nil
}

/**
 * block 找到索引文件，更新字典 start 文件所以位置
 * @param  {[type]} m *mpso),id 当前文章,id previouid 文章上一个Id,
 * count 关键词次数,关键词开始位置，关键词结束位置
 * @return {[type]}   [description]
 **/

func (m *mpso) Set(id, previouid, count, strid, endid int32) (ns int32, block int32, err error) {

	file := fmt.Sprintf("%s/%d.bin", m.path, m.idx)

	fi, err := os.Stat(file)
	if err != nil {
		return 0, 0, err
	}

	m.mutex.Lock()
	n, err := m.w.Seek(0, os.SEEK_END)

	//fmt.Println("fi size", fi.Size())

	if fi.Size() > MpBlockSize {

		m.w.Close()
		m.idx += 1

		file = fmt.Sprintf("%s/%d.bin", m.path, m.idx)
		m.w, err = os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0664)

		if err != nil {
			m.idx -= 1
			return 0, 0, err
		}
		n = 0
	}

	var sid, previousid, bcount, bufstr, bufend []byte //= make([]byte, 4)

	sid, previousid, bcount, bufstr, bufend = make([]byte, 4), make([]byte, 4), make([]byte, 4), make([]byte, 4), make([]byte, 4)

	binary.BigEndian.PutUint32(sid, uint32(id))
	binary.BigEndian.PutUint32(previousid, uint32(previouid))
	binary.BigEndian.PutUint32(bcount, uint32(count))
	binary.BigEndian.PutUint32(bufstr, uint32(strid))
	binary.BigEndian.PutUint32(bufend, uint32(endid))

	bf := bytes.NewBuffer([]byte{})

	bf.Write(sid)
	bf.Write(previousid)
	bf.Write(bcount)
	bf.Write(bufstr)
	bf.Write(bufend)

	//m.w.Seek(int64(start), 0)
	//写入文件
	//fmt.Println("byt:", bf.Bytes())
	//n, err := m.w.Write(bf.Bytes())
	_, err = m.w.WriteAt(bf.Bytes(), n)

	defer m.mutex.Unlock()

	if err != nil {
		return 0, m.idx, err
	}
	return int32(n), m.idx, nil
}

/**
 * 更新map 的索引关系
 * key 关键词，start 开始位置,idx 搜索位置
 **/
func (m *mpso) Update(key string, start, idx int32) error {

	//获取路径词典索引路径
	//更新v
	//
	m.dicmap[key] = fmt.Sprintf("%d:%d", idx, start)

	return nil
	/**
	if m.dicmapmax < 100 {
		return nil
	}

	dicidxjson, err := json.Marshal(m.dicmap)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	m.w.Seek(int64(Hder), 0)

	_, err = m.w.Write(dicidxjson)
	m.mutex.Unlock()
	if err != nil {
		return err
	}
	m.dicmapmax = 0
	return nil

	**/
}

func (m *mpso) Closemap() error {

	dicidxjson, err := json.Marshal(m.dicmap)
	if err != nil {
		return err
	}

	dictmpfileName := fmt.Sprintf("%s/%s", m.path, mpdict)

	fi, err := os.OpenFile(dictmpfileName, os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		//fmt.Println("dictmpfileName", dictmpfileName, err)
		return err
	}

	fi.Seek(int64(Hder), 0)

	bf := bytes.NewBuffer([]byte{})

	bf.Write(dicidxjson)

	_, err = fi.Write(bf.Bytes())

	if err != nil {
		//fmt.Println(err)
		return err
	}

	return nil
}

//查询关键词对应的文章id
func (m *mpso) Search(key string) (block int32, id int32, err error) {
	v, ok := m.dicmap[key]
	if ok == false {
		return 0, 0, errors.New("error key not ")
	}

	vr := strings.Split(v, ":")
	a, _ := strconv.Atoi(vr[0])
	b, _ := strconv.Atoi(vr[1])

	return int32(a), int32(b), nil
}

//id 是起始位置
func (m *mpso) Limit(id, block, limit int32) ([100]int32, error) {
	var j [100]int32
	var i int32

	for i = 0; i < limit; i++ {
		n1, n2, _, _, _, err := m.Get(block, id)

		//log.Infof("n1 %d,n2 %d", n1, n2)

		if err != nil {
			return j, err
		}

		if n2 <= 0 {
			break
		}

		id = n2
		j[i] = n1
	}
	return j, nil
}
