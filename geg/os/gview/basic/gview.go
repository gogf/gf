package main


import (
    "fmt"
    "gitee.com/johng/gf/g"
)

func main() {
    v      := g.View()
    b, err := v.Parse("gview.tpl", map[string]interface{} {
        "k" : "v",
    })
    fmt.Println(err)
    fmt.Println(string(b))
}