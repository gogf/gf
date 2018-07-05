package main

import (
    "gitee.com/johng/gf/g/util/gregx"
    "fmt"
)

func main() {
    fmt.Println(gregx.MatchString(`[-/]`, "-"))
}
