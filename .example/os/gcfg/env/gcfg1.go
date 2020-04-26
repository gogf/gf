package main

import (
	"fmt"

	"github.com/gogf/gf/os/genv"

	"github.com/gogf/gf/frame/g"
)

// 使用第二个参数指定读取的配置文件
func main() {
	genv.Set("PPROF_ENABLE", "true")
	genv.Set("REDIS_HOST", "localhost")
	genv.Set("REDIS_PORT", "6378")
	c := g.Config()
	fmt.Println(c.GetBool("server.PProfEnabled"))
	fmt.Println(c.GetArray("redis-cache"))

	genv.Remove("PPROF_ENABLE")
	genv.Remove("REDIS_HOST")
	genv.Remove("REDIS_PORT")
	fmt.Println(c.GetBool("server.PProfEnabled"))
	fmt.Println(c.GetArray("redis-cache"))
}
