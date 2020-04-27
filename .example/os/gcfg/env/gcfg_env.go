package main

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
)

// 使用第二个参数指定读取的配置文件
func main() {
	c := g.Config()
	fmt.Println(c.GetBool("server.PProfEnabled"))
	fmt.Println(c.GetArray("redis-cache"))
}
