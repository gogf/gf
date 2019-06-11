package main

import (
<<<<<<< HEAD
    "fmt"
    "gitee.com/johng/gf/g"
)

func main() {
    fmt.Println(g.Config().Get("serverpath"))
}

=======
	"fmt"
	"github.com/gogf/gf/g"
)

// 使用GetVar获取动态变量
func main() {
	fmt.Println(g.Config().GetVar("memcache.0").String())
}
>>>>>>> upstream/master
