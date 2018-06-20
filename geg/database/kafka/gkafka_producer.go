package main

import (
    "gitee.com/johng/gf/g/database/gkafka"
    "fmt"
)

func main () {
    config        := gkafka.NewConfig()
    config.Servers = "localhost:9092"
    config.Topics  = "test"

    client := gkafka.NewClient(config)
    defer client.Close()
    err := client.SyncSend(&gkafka.Message{Value: []byte("1")})
    fmt.Println(err)
}
