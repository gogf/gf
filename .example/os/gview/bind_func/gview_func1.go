package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/os/gview"
)

// 用于测试的内置函数
func funcTest() string {
	return "test content"
}

func main() {
	// 解析模板的时候传递模板函数映射Map，仅会在当前模板解析生效
	parsed, err := g.View().ParseContent(`call build-in function test: {{test}}`, nil, gview.FuncMap{
		"test": funcTest,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(string(parsed))
}
