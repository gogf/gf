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
    if err != nil {
        panic(err)
    }
    fmt.Println(string(b))
}