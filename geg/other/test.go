package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gjson"
)

func main() {
    j, _ := gjson.Load("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame/config.json")
    c, _ := j.ToXmlIndent("config")
    fmt.Println(string(c))
}