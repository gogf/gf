package main

import (
<<<<<<< HEAD
    "fmt"
    "gitee.com/johng/gf/g/os/gcfg"
)

func main() {
    c := gcfg.New("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/os/gcfg")
    fmt.Println(c.GetArray("memcache"))
}

=======
	"fmt"
	"github.com/gogf/gf/g"
)

// 使用默认的config.toml配置文件读取配置
func main() {
	c := g.Config()
	fmt.Println(c.GetArray("memcache"))
}
>>>>>>> upstream/master
