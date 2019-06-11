package main

import (
<<<<<<< HEAD
    "fmt"
    "gitee.com/johng/gf/g/os/gcfg"
)

func main() {
    c              := gcfg.New("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/os/gcfg")
    redisConfig    := c.GetArray("redis-cache", "redis.yml")
    memcacheConfig := c.GetArray("", "memcache.yml")
    fmt.Println(redisConfig)
    fmt.Println(memcacheConfig)
}

=======
	"fmt"
	"github.com/gogf/gf/g"
)

// 使用第二个参数指定读取的配置文件
func main() {
	c := g.Config()
	redisConfig := c.GetArray("redis-cache", "redis.toml")
	memConfig := c.GetArray("", "memcache.yml")
	fmt.Println(redisConfig)
	fmt.Println(memConfig)
}
>>>>>>> upstream/master
