package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
)

// 使用GetVar获取动态变量
func main() {
	fmt.Println(g.Config().GetVar("memcache.0").String())
}
