package main

import (
    "gitee.com/johng/gf/g/database/gkafka"
    "fmt"
)

func main () {
    config        := gkafka.NewConfig()
    config.GroupId = "group_2"
    config.Servers = "localhost:9092"
    config.Topics  = "abc"

    client := gkafka.New(config)
    defer client.Close()

    for {
       msg, err := client.Receive()
       fmt.Printf("value: %s, err: %v", string(msg.Value), err)
    }
}
