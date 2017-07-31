package main

import (
    "fmt"
)

type LogEntry struct {
    id               int64        // 唯一ID
    act              string       // 操作
    key              string
    value            string
}

func main() {
    var a LogEntry
    a.id = 1
    fmt.Println(a)
}