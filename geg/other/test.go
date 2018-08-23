package main

import (
    "gitee.com/johng/gf/g/util/gregex"
    "fmt"
)


func main() {
    s := `<a href="baidu.com">百度</a>"`
    gregex.ReplaceStringFunc(`href="(.+?)"`, s, func(s string) string {
        fmt.Println(s)
        return s
    })
}