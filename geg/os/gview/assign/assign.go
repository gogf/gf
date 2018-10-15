package main


import (
    "fmt"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gview"
)

func main() {
    g.View().Assigns(gview.Params{
        "k1" : "v1",
        "k2" : "v2",
    })
    b, err := g.View().ParseContent(`{{.k1}} - {{.k2}}`, nil)
    fmt.Println(err)
    fmt.Println(string(b))
}