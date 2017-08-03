package main

import (
    "fmt"
    "g/util/grand"
)

type LogEntry struct {
    id               int64        // 唯一ID
    act              string       // 操作
    key              string
    value            string
}

func main() {
    for i := 0; i < 10; i++ {
        fmt.Println(grand.Rand(0, 1))
    }
}