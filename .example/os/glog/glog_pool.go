package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// 测试删除日志文件是否会重建日志文件
func main() {
	path := "/Users/john/Temp/test"
	g.Log().SetPath(path)
	for {
		g.Log().Print(gtime.Now().String())
		time.Sleep(time.Second)
	}
}
