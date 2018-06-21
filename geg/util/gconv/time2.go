package main

import (
    "gitee.com/johng/gf/g/util/gconv"
    "fmt"
)

func main() {
    fmt.Println(gconv.Time("2018-06-07").String())

    fmt.Println(gconv.Time("2018-06-07 13:01:02").String())

    fmt.Println(gconv.Time("2018-06-07 13:01:02.096").String())


}