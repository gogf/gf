package main

import (
    "gitee.com/johng/gf/g/database/gkafka"
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
)

func main () {
    config        := gkafka.NewConfig()
    config.GroupId = "group_1"
    config.Servers = "localhost:9092"
    config.Topics  = "test"
    config.AutoMarkOffset = false

    client := gkafka.NewClient(config)
    defer client.Close()

    for {
       msg, err := client.Receive()
       fmt.Printf("%s value: %s, err: %v\n", gtime.Datetime(), string(msg.Value), err)
       msg.MarkOffset()
    }
}
