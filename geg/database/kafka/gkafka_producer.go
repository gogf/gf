package main

import (
    "gitee.com/johng/gf/g/database/gkafka"
    "fmt"
)

func main () {
    client := gkafka.New(gkafka.Config{
        Servers : "localhost:9092",
        Topics  : "abc",
    })
    defer client.Close()
    err := client.SyncSend(&gkafka.Message{Value: []byte("111")})
    fmt.Println(err)
}
