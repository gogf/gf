package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
)

// 使用g.Config方法获取配置管理对象，并指定默认的配置文件名称
func main() {
	fmt.Println(g.Config().Get("viewpath"))
}
