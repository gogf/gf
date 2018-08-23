package main


import (
    "fmt"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gview"
)

// 用于测试的内置函数
func funcTest() string {
    return "test"
}

func main() {
    view   := g.View()
    b, err := view.Parse("gview.tpl", nil, gview.FuncMap{
        "test" : funcTest,
    })
    fmt.Println(err)
    fmt.Println(string(b))
}