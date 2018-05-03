package main


import (
    "fmt"
    "gitee.com/johng/gf/g"
)

func main() {
    v      := g.View()
    b, err := v.Parse("test.tpl", nil)
    fmt.Println(err)
    fmt.Println(b)
}