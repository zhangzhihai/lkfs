package lkfs

import (
	"fmt"
	"github.com/huichen/sego"
	"strconv"
	"strings"
	"testing"
)

var fi *mpso

func TestNewmap(t *testing.T) {
	var err error
	path := "d:/webserver/net/golang/src/ulucu.github.com/lkfs/data"

	fi, err = Newmap(path)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("...")

}

func TestMapSet(t *testing.T) {

	var n int32 = 0
	var b int32 = 0
	var err error

	for j := 0; j < 100; j++ {
		n, b, err = fi.Set(int32(10+j), n, 3, 1, 5)
		if err != nil {
			t.Error(err)
		}
		//fmt.Println("..n:", n)
	}

	fmt.Println(n, b)
}

/**
func TestMapGet(t *testing.T) {

	n1, n2, n3, n4, n5, err := fi.Get(1, 6891160)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(n1, n2, n3, n4, n5)
}

func TestUpdateso(t *testing.T) {
	err := fi.Update("hello", 6891160, 1)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("testsoupdate")
}

func TestGetso(t *testing.T) {

	i, j, err := fi.Search("hello")
	if err != nil {
		t.Error(err)
	}
	fmt.Println("testsoGet")
	fmt.Printf("i,j %d %d", i, j)
}

func TestLimit(t *testing.T) {
	//100起始位置

	i, err := fi.Limit(1980, 0, 10)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("i,j %v", i)
}
**/

//完整写入方式
func TestSegmen(t *testing.T) {
	var segmenter sego.Segmenter
	segmenter.LoadDictionary("d:/webserver/net/golang/src/ulucu.github.com/lkfs/data/dictionary.txt")

	// 分词
	text := []byte(`近年结识了一位警察朋友，好枪法。不单单在射击场上百发百中，更在解救人质的现场，次次百步穿杨。当然了，这个“杨”不是杨树的杨，而是匪徒的代称。
　　我向他请教射击的要领。他说，很简单，就是极端的平静。我说这个要领所有打枪的人都知道，可是做不到。他说，记住，你要像烟灰一样松散。只有放松，全部潜在的能量才会释放出来，协同你达到完美。
　　他的话我似懂非懂，但从此我开始注意以前忽略了的烟灰。烟灰，尤其是那些优质香烟燃烧后的烟灰，非常松散，几乎没有重量和形状，真一个大象无形。它们懒洋洋地趴在那里，好像在冬眠。其实，在烟灰的内部，栖息着高度警觉和机敏的鸟群，任何一阵微风掠过，哪怕只是极轻微的叹息，它们都会不失时机地腾空而起驭风而行。它们的力量来自放松，来自一种飘扬的本能。
　　松散的反面是紧张。几乎每个人都有过由于紧张而惨败的经历。比如，考试的时候，全身肌肉僵直，心跳得好像无数个小炸弹在身体的深浅部位依次爆破。手指发抖头冒虚汗，原本记得滚瓜烂熟的知识，改头换面潜藏起来，原本泾渭分明的答案变得似是而非，泥鳅一样滑走……面试的时候，要么扭扭捏捏不够大方，无法表现自己的真实实力，要么口若悬河躁动不安，拿捏不准问题的实质，只得用不停的述说掩饰自己的紧张，适得其反……相信每个人都储存了一大堆这类不堪回首的往事。在最危急的时刻能保持极端的放松，不是一种技术，而是一种修养，是一种长期潜移默化修炼提升的结果。我们常说，某人胜就胜在心理上，或是说某人败就败在心理上。这其中的差池不是指在理性上，而是这种心灵张弛的韧性上。
　　没事的时候看看烟灰吧。他们曾经是火焰，燃烧过，沸腾过，但它们此刻安静了。它们毫不张扬地聚精会神地等待着下一次的乘风而起，携带着全部的能量，抵达阳光能到的任何地方`)
	segments := segmenter.Segment(text)

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	out := sego.SegmentsToSlice(segments, false)

	var err error
	path := "d:/webserver/net/golang/src/ulucu.github.com/lkfs"

	//把txt 写入存储，返回的id
	news, err := New(path)
	if err != nil {
		t.Error(err)
	}

	id, err := news.Set(text)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(id)

	//把分词写入索引
	for _, v := range out {
		if len(v) > 3 {
			//将分词写入词库
			//n 上一篇文章位置
			x, ok := fi.dicmap[v]
			if ok {

				vr := strings.Split(x, ":")

				n, _ := strconv.Atoi(vr[1]) //上一篇文章位置
				//fmt.Println("n", n)

				nxxx, bxxx, err := fi.Set(id, int32(n), 1, 0, 0)
				if err != nil {
					t.Error(err)
				}

				err = fi.Update(v, nxxx, bxxx)
				if err != nil {
					t.Error(err)
				}

				//fmt.Println("nxxx,bxxx", nxxx, bxxx)
			}

		}
	}
	err = fi.Closemap()
	if err != nil {
		fmt.Println("close fail")
		t.Error(err)
	}
	fmt.Println("close ok")

}

////完整读取方式
func TestSoreader(t *testing.T) {
	//key := []byte("百发百中")
	//搜索 安静,i bolck,j位置
	i, j, err := fi.Search("枪法")
	if err != nil {
		t.Error(err)
	}
	//i文章的索引
	fmt.Println(i, j)
	//j 关键索引的位置
	idx, err := fi.Limit(j, i, 10)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("idx .....", idx)

	n1, n2, n3, n4, n5, err := fi.Get(i, j)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(n1, n2, n3, n4, n5)

	path := "d:/webserver/net/golang/src/ulucu.github.com/lkfs"

	//把txt 写入存储，返回的id
	news, err := New(path)
	//从指定位置取文章内容
	body, err := news.Get(n1)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(body))

}
