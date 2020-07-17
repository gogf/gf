package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
)

// 用于测试的带参数的内置函数
func funcHello(name string) string {
	return fmt.Sprintf(`Hello %s`, name)
}

func main() {
	// 绑定全局的模板函数
	g.View().BindFunc("hello", funcHello)

	// 普通方式传参
	parsed1, err := g.View().ParseContent(`{{hello "GoFrame"}}`, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(parsed1))

	// 通过管道传参
	parsed2, err := g.View().ParseContent(`{{"GoFrame" | hello}}`, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(parsed2))
}
