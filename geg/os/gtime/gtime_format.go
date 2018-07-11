package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    formats := []string{
        "Y-m-d H:i:s.u",
        "D M d H:i:s T O Y",
        // 可以使用转义字符转移有意义的格式字母
        "\\T\\i\\m\\e \\i\\s: h:i:s a",
        // format格式不支持标准库格式混合，相互隔离
        "2006-01-02T15:04:05.000000000Z07:00",
    }
    t := gtime.Now()
    for _, f := range formats {
        fmt.Println(f)
        fmt.Println(t.Format(f))
        fmt.Println()
    }
}
