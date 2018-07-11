package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
)

func main() {
    formats := []string{
        "Y-m-d H:i:s",
        "2006-01-02T15:04:05Z07:00",
        "2006-01-02T15:04:05.999999999Z07:00",
    }
    t := gtime.Now()
    for _, f := range formats {
        fmt.Println(f)
        fmt.Println(t.Format(f))
        fmt.Println()
    }
}
