package lkfs

import (
//"fmt"
//"testing"
)

var fi *filesystem
var block, start, size int32
var body []byte

/**
func init() {
	block = 0
	start = 0
	size = 39
	body = []byte("hello world 这里有多少个字符串")
}

func TestNewfs(t *testing.T) {
	var err error
	path := "E:/net/golang/src/ulucu.github.com/lkfs"
	fi, err = Newfs(path)
	if err != nil {
		t.Log(err)
	}

	//fmt.Println(pt)
}

func TestRead(t *testing.T) {
	for i := int32(0); i < 10; i++ {
		start = i * size
		b, err := fi.Read(block, start, size)
		if err != nil {
			t.Log(err)
		}

		fmt.Println("i", string(b))
	}
}

func TestFWrite(t *testing.T) {

	for i := int32(0); i < 10; i++ {
		start = i * size
		bl, err := fi.Write(block, start, size, body)

		fmt.Println("len", len(body))
		if err != nil {
			t.Log(err)
		}
		fmt.Println(bl)
	}

}
**/
