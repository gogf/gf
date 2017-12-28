package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gregx"
)


func main() {
    fmt.Println(gregx.IsMatchString(`^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)+$`, "joh-n_cn@johng.cn"))
}