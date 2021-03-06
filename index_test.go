package lkfs

import (
	"fmt"
	"testing"
)

//go test -bench=".*" -file index_test.go index.go
// go test . -bench=".*"
//测试某个方法 go test -run='Test_xxx'
var bi *FileIdx

func TestInit(t *testing.T) {

	path := "d:/webserver/net/golang/src/ulucu.github.com/lkfs"
	pt, err := Init(path)
	if err != nil {
		t.Log(err)
	}

	fmt.Println(pt)
}

func TestNewIndex(t *testing.T) {
	path := "d:/webserver/net/golang/src/ulucu.github.com/lkfs"
	pt, err := NewFildIdx(path)
	if err != nil {
		t.Error(err)
		fmt.Println(err.Error())
	}
	bi = pt
	fmt.Println("test", pt)
}

func TestUuid(t *testing.T) {

	uid, err := bi.Uuid()

	if err != nil {
		t.Log(err)
	}
	fmt.Println(uid)

}

func TestWrite(t *testing.T) {

	var i int32
	for i = 0; i < 50; i++ {

		uid, err := bi.Uuid()
		fmt.Println("uid:", uid)

		err = bi.Write(uid*10000, 0, i*int32(50), 1024)
		if err != nil {
			t.Log(err)
		}
		fmt.Println(uid)
	}

}

func TestReader(t *testing.T) {
	b, s, l, err := bi.Reader(20)

	if err != nil {
		t.Log(err)
	}

	fmt.Println("b", b, "s", s, "l", l)
}
