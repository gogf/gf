package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gtime"
	"time"
)

// 配置文件热更新示例
func main() {
	c := g.Config()
	// 每隔1秒打印当前配置项值，用户可手动在外部修改文件内容，gcfg读取到的配置项值会即时得到更新
	gtime.SetInterval(time.Second, func() bool {
		fmt.Println(c.Get("viewpath"))
		return true
	})

	select {}
}
