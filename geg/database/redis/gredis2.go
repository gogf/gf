package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gconv"
)

// 使用框架封装的g.Redis()方法获得redis操作对象单例，不需要开发者显示调用Close方法
func main() {
	g.Redis().Do("SET", "k", "v")
	v, _ := g.Redis().Do("GET", "k")
	fmt.Println(gconv.String(v))
}
