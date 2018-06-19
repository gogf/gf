package main

import (
    "gitee.com/johng/gf/g/database/gkafka"
    "fmt"
)

func main () {
    config        := gkafka.NewConfig()
    config.Servers = "localhost:9092"
    config.Topics  = "abc"

    client := gkafka.New(config)
    defer client.Close()
    err := client.SyncSend(&gkafka.Message{Value: []byte("111")})
    fmt.Println(err)
}
