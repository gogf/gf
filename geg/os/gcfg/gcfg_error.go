package main

import (
	"fmt"
	"github.com/gogf/gf/g"
)

// 演示在找不到配置文件时的错误提示
func main() {
	fmt.Println(g.Config("none-exist-config.toml").Get("none"))
}
