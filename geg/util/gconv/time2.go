package main

import (
    "gitee.com/johng/gf/g/util/gconv"
    "fmt"
)

func main() {
    fmt.Println(gconv.Time("2018-06-07").Date())
    fmt.Println(gconv.Time("2018-06-07").Clock())
    fmt.Println(gconv.Time("2018-06-07 13:01:02").Date())
    fmt.Println(gconv.Time("2018-06-07 13:01:02").Clock())
}