package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gvalid"
)

func main() {
    data  := map[string]interface{} {
        "id"   : "1",
    }
    rules := map[string]string {
        "id"   : "required",
        "name" : "length:4,16",
    }
    m := gvalid.CheckMap(data, rules)
    fmt.Println(m)
}
