package main

import (
	"time"

	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/container/garray"
)

func main() {
	// 创建对象
	var a = garray.NewStrArray()
	// push 100个 字符串进去
	for i := 1; i <= 10; i++ {
		a.PushLeft("a_" + gconv.String(i))
	}
	// 死循环 Pop 取出
	for {
		glog.Printf("a.Len() ---> %d", a.Len())
		if a.Len() == 0 {
			break
		}
		if v, isFound := a.PopRight(); isFound {
			glog.Printf("Pop -----> %s", v)
		}
		time.Sleep(50 * time.Millisecond)
	}
}
