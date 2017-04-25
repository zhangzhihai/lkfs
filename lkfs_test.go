package lkfs

import (
	//"bytes"
	//"encoding/binary"
	"fmt"
	//"io/ioutil"
	//"os"
	"testing"
)

var news *lkfs

func TestNew(t *testing.T) {
	var err error
	path := "d:/webserver/net/golang/src/ulucu.github.com/lkfs"

	/**
		//pt := "d:/webserver/net/golang/src/ulucu.github.com/lkfs/test.txt"
		pt := "d:/webserver/net/golang/src/ulucu.github.com/lkfs/bin.idx"
		dat, err := ioutil.ReadFile(pt)
		fmt.Println(dat)
		//fi, err := os.OpenFile(pt, os.O_RDWR|os.O_CREATE, 0664)
		fi, err := os.Open(pt)


			buuid := make([]byte, 4)
			binary.BigEndian.PutUint32(buuid, uint32(5010000))

			bf := bytes.NewBuffer(make([]byte, 0))
			bf.Write(buuid)

			fmt.Println(bf.Bytes())

			n, err := fi.Write(bf.Bytes())


		fmt.Println("=========================")

		buf := make([]byte, 4)
		n, err := fi.Read(buf)
		fmt.Println("n....", n)
		fmt.Println("out buf ", buf)
		//bbbb := bytes.NewReader(buf)
		var pi int32
		//err = binary.Read(bbbb, binary.LittleEndian, &pi)
		binary.Read(bytes.NewBuffer(buf[0:4]), binary.BigEndian, &pi)
		fmt.Println(pi)

		fi.Close()
	**/
	news, err = New(path)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(news)
}

func TestSet(t *testing.T) {
	//rs := "D:/webserver/net/word/OpenResty-Best-Practices-20160812.pdf"
	//body, err := ioutil.ReadFile(rs)
	body := []byte(`golang判断文件或文件夹是否存在的方法为使用os.Stat()函数返回的错误值进行判断:`)
	id, err := news.Set(body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("...")
	t.Log(id)

}

func TestGet(t *testing.T) {

	body, err := news.Get(35)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(body))

}
